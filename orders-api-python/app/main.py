from fastapi import FastAPI, Request
from fastapi.responses import JSONResponse
import time

from app.api.routes import router as orders_router
from app.core.config import get_settings
from app.database.mock_db import add_sample_orders


def create_app() -> FastAPI:
    """
    Create and configure the FastAPI application
    """
    settings = get_settings()

    app = FastAPI(
        title=settings.app_name,
        description="A simple REST API for managing orders",
        version=settings.version,
        debug=settings.debug,
        docs_url="/docs",
        redoc_url="/redoc",
    )

    @app.middleware("http")
    async def add_process_time_header(request: Request, call_next):
        """Add X-Process-Time header to responses"""
        start_time = time.time()
        response = await call_next(request)
        process_time = time.time() - start_time
        response.headers["X-Process-Time"] = str(process_time)
        return response

    @app.get("/status", tags=["health"])
    async def health_check():
        """Health check endpoint"""
        return {
            "status": "ok",
            "name": settings.app_name,
            "version": settings.version
        }

    # Include API routes
    app.include_router(orders_router, prefix=settings.api_prefix)

    # Add sample data for demonstration
    @app.on_event("startup")
    def startup_event():
        """Add sample data on startup"""
        add_sample_orders()

    return app


app = create_app()


if __name__ == "__main__":
    # This is used when running locally with python main.py
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8000)
