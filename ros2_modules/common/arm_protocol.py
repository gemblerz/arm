"""Shared JSON protocol helpers for ROS2 arm modules."""

from __future__ import annotations

from dataclasses import dataclass
from datetime import datetime, timezone
from typing import Dict


def utc_timestamp() -> str:
    """Return current UTC timestamp in ISO-8601 format."""
    return datetime.now(timezone.utc).isoformat()


def clamp(value: float, minimum: float, maximum: float) -> float:
    """Clamp a float to a closed interval."""
    return max(minimum, min(maximum, value))


@dataclass(frozen=True)
class CartesianTarget:
    """End-effector Cartesian target."""

    x: float
    y: float
    z: float


@dataclass(frozen=True)
class JointTargets:
    """Named joint targets in degrees."""

    joints: Dict[str, float]
