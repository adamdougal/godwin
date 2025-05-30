from functools import lru_cache
from typing import Dict, Any, Optional

from pydantic_settings import BaseSettings, SettingsConfigDict
from dotenv import find_dotenv


class Settings(BaseSettings):
    """
    Application settings.

    These settings can be configured using environment variables or a .env file.
    """
    app_name: str = "Orders API"
    debug: bool = False
    version: str = "0.1.0"
    api_prefix: str = "/api/v1"

    # Configure using environment variables with prefix "ORDERS_"
    model_config = SettingsConfigDict(
        env_prefix="orders_",
        env_file=find_dotenv(".env"),
        env_file_encoding="utf-8",
        case_sensitive=False,
    )


@lru_cache()
def get_settings() -> Settings:
    """
    Get the application settings.

    Uses lru_cache to avoid re-reading the environment each time settings are accessed.
    """
    return Settings()
