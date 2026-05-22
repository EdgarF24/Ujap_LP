import uuid
from datetime import datetime
from typing import Optional

from pydantic import BaseModel, EmailStr, Field

from app.models.user import UserRole


# ---------------------------------------------------------------------------
# Request schemas
# ---------------------------------------------------------------------------


class UserCreate(BaseModel):
    email: EmailStr
    password: str = Field(..., min_length=6, description="Minimum 6 characters")
    full_name: str = Field(..., min_length=1, max_length=255)


class UserLogin(BaseModel):
    email: EmailStr
    password: str


class UserUpdate(BaseModel):
    full_name: Optional[str] = Field(None, min_length=1, max_length=255)
    role: Optional[UserRole] = None
    is_active: Optional[bool] = None


# ---------------------------------------------------------------------------
# Response schemas
# ---------------------------------------------------------------------------


class UserResponse(BaseModel):
    id: uuid.UUID
    email: EmailStr
    full_name: str
    role: UserRole
    is_active: bool
    created_at: datetime

    model_config = {"from_attributes": True}


# ---------------------------------------------------------------------------
# Token schemas
# ---------------------------------------------------------------------------


class Token(BaseModel):
    access_token: str
    token_type: str = "bearer"


class TokenData(BaseModel):
    email: str
    role: UserRole


# ---------------------------------------------------------------------------
# Internal validation schema
# ---------------------------------------------------------------------------


class TokenValidateRequest(BaseModel):
    token: str


class TokenValidateResponse(BaseModel):
    valid: bool
    user_id: Optional[str] = None
    email: Optional[str] = None
    role: Optional[str] = None
