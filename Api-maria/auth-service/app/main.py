from contextlib import asynccontextmanager

from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware

from app.db.database import Base, engine
from app.routers.auth import router as auth_router


@asynccontextmanager
async def lifespan(app: FastAPI):
    """Create database tables on startup (Alembic handles migrations in prod)."""
    Base.metadata.create_all(bind=engine)
    yield


app = FastAPI(
    title="Co-Working Auth Service",
    description="JWT-based authentication microservice for the co-working platform.",
    version="1.0.0",
    lifespan=lifespan,
)

# ---------------------------------------------------------------------------
# CORS – allow all origins for development
# ---------------------------------------------------------------------------
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# ---------------------------------------------------------------------------
# Routers
# ---------------------------------------------------------------------------
app.include_router(auth_router, prefix="/auth")


# ---------------------------------------------------------------------------
# Health check
# ---------------------------------------------------------------------------


@app.get("/health", tags=["health"], summary="Service health check")
def health_check() -> dict:
    """Return a simple liveness probe response."""
    return {"status": "healthy", "service": "auth-service"}
