import uuid
from typing import List

from fastapi import APIRouter, Depends, HTTPException, Query, status
from sqlalchemy.orm import Session

from app.core.dependencies import get_current_user, require_admin
from app.core.security import (
    create_access_token,
    decode_token,
    get_password_hash,
    verify_password,
)
from app.db.database import get_db
from app.models.user import User, UserRole
from app.schemas.user import (
    Token,
    TokenValidateRequest,
    TokenValidateResponse,
    UserCreate,
    UserLogin,
    UserResponse,
    UserUpdate,
)

router = APIRouter(tags=["auth"])

ADMIN_EMAIL = "admin@coworking.com"


# ---------------------------------------------------------------------------
# POST /auth/register
# ---------------------------------------------------------------------------


@router.post(
    "/register",
    response_model=UserResponse,
    status_code=status.HTTP_201_CREATED,
    summary="Register a new user",
)
def register(user_in: UserCreate, db: Session = Depends(get_db)) -> UserResponse:
    """Create a new user account.

    * The first registration with ``admin@coworking.com`` automatically
      receives the **ADMIN** role.
    * All other accounts receive the **USER** role by default.
    """
    existing = db.query(User).filter(User.email == user_in.email).first()
    if existing:
        raise HTTPException(
            status_code=status.HTTP_409_CONFLICT,
            detail="An account with this email already exists",
        )

    role = UserRole.ADMIN if user_in.email == ADMIN_EMAIL else UserRole.USER

    new_user = User(
        email=user_in.email,
        hashed_password=get_password_hash(user_in.password),
        full_name=user_in.full_name,
        role=role,
    )
    db.add(new_user)
    db.commit()
    db.refresh(new_user)
    return new_user


# ---------------------------------------------------------------------------
# POST /auth/login
# ---------------------------------------------------------------------------


@router.post(
    "/login",
    response_model=Token,
    summary="Authenticate and receive a JWT",
)
def login(credentials: UserLogin, db: Session = Depends(get_db)) -> Token:
    """Verify email/password and return a Bearer JWT access token."""
    user: User | None = (
        db.query(User).filter(User.email == credentials.email).first()
    )
    if user is None or not verify_password(credentials.password, user.hashed_password):
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail="Incorrect email or password",
            headers={"WWW-Authenticate": "Bearer"},
        )

    if not user.is_active:
        raise HTTPException(
            status_code=status.HTTP_403_FORBIDDEN,
            detail="User account is inactive",
        )

    access_token = create_access_token(
        data={"sub": user.email, "role": user.role.value, "user_id": str(user.id)}
    )
    return Token(access_token=access_token, token_type="bearer")


# ---------------------------------------------------------------------------
# GET /auth/me
# ---------------------------------------------------------------------------


@router.get(
    "/me",
    response_model=UserResponse,
    summary="Get the current authenticated user profile",
)
def get_me(current_user: User = Depends(get_current_user)) -> UserResponse:
    """Return the profile of the currently authenticated user."""
    return current_user


# ---------------------------------------------------------------------------
# POST /auth/validate  (internal endpoint)
# ---------------------------------------------------------------------------


@router.post(
    "/validate",
    response_model=TokenValidateResponse,
    summary="Internal: validate a JWT token",
)
def validate_token(body: TokenValidateRequest) -> TokenValidateResponse:
    """Validate a JWT token and return decoded identity claims.

    This endpoint is intended for other microservices to verify tokens
    without sharing the JWT secret directly.
    """
    payload = decode_token(body.token)
    if payload is None:
        return TokenValidateResponse(valid=False)

    email: str | None = payload.get("sub")
    role: str | None = payload.get("role")
    user_id: str | None = payload.get("user_id")

    if email is None:
        return TokenValidateResponse(valid=False)

    return TokenValidateResponse(
        valid=True,
        user_id=user_id,
        email=email,
        role=role,
    )


# ---------------------------------------------------------------------------
# GET /auth/users  (admin)
# ---------------------------------------------------------------------------


@router.get(
    "/users",
    response_model=List[UserResponse],
    summary="Admin: list all users",
)
def list_users(
    skip: int = Query(0, ge=0, description="Number of records to skip"),
    limit: int = Query(20, ge=1, le=100, description="Maximum records to return"),
    db: Session = Depends(get_db),
    _admin: User = Depends(require_admin),
) -> List[UserResponse]:
    """Return a paginated list of all registered users. Requires ADMIN role."""
    users = db.query(User).offset(skip).limit(limit).all()
    return users


# ---------------------------------------------------------------------------
# GET /auth/users/{user_id}  (admin)
# ---------------------------------------------------------------------------


@router.get(
    "/users/{user_id}",
    response_model=UserResponse,
    summary="Admin: get a specific user",
)
def get_user(
    user_id: uuid.UUID,
    db: Session = Depends(get_db),
    _admin: User = Depends(require_admin),
) -> UserResponse:
    """Return a single user by UUID. Requires ADMIN role."""
    user: User | None = db.query(User).filter(User.id == user_id).first()
    if user is None:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail="User not found",
        )
    return user


# ---------------------------------------------------------------------------
# PUT /auth/users/{user_id}  (admin)
# ---------------------------------------------------------------------------


@router.put(
    "/users/{user_id}",
    response_model=UserResponse,
    summary="Admin: update a user",
)
def update_user(
    user_id: uuid.UUID,
    user_in: UserUpdate,
    db: Session = Depends(get_db),
    _admin: User = Depends(require_admin),
) -> UserResponse:
    """Update a user's full_name, role, or is_active flag. Requires ADMIN role."""
    user: User | None = db.query(User).filter(User.id == user_id).first()
    if user is None:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail="User not found",
        )

    update_data = user_in.model_dump(exclude_unset=True)
    for field, value in update_data.items():
        setattr(user, field, value)

    db.commit()
    db.refresh(user)
    return user


# ---------------------------------------------------------------------------
# DELETE /auth/users/{user_id}  (admin – soft delete)
# ---------------------------------------------------------------------------


@router.delete(
    "/users/{user_id}",
    status_code=status.HTTP_200_OK,
    summary="Admin: soft-delete a user",
)
def delete_user(
    user_id: uuid.UUID,
    db: Session = Depends(get_db),
    _admin: User = Depends(require_admin),
) -> dict:
    """Soft-delete a user by setting ``is_active = False``. Requires ADMIN role."""
    user: User | None = db.query(User).filter(User.id == user_id).first()
    if user is None:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail="User not found",
        )

    user.is_active = False
    db.commit()
    return {"detail": "User deactivated successfully"}
