"""Robot controller kinematics helpers."""

from __future__ import annotations

import math

from ros2_modules.common.arm_protocol import CartesianTarget, JointTargets, clamp


class RobotControllerKinematics:
    """Small inverse-kinematics helper for ROS2 joint targeting."""

    def __init__(self, max_reach_xy: float = 0.8) -> None:
        self.max_reach_xy = max_reach_xy

    def solve(self, target: CartesianTarget) -> JointTargets:
        """Convert Cartesian target to simple 3-joint targets."""
        radius = math.hypot(target.x, target.y)
        limited_radius = clamp(radius, 0.0, self.max_reach_xy)
        scale = limited_radius / radius if radius > 0 else 0.0
        clamped_x = target.x * scale
        clamped_y = target.y * scale

        base = math.degrees(math.atan2(clamped_y, clamped_x)) if radius > 0 else 0.0
        shoulder = clamp((limited_radius / self.max_reach_xy) * 90.0, 0.0, 90.0)
        elbow = clamp((target.z + 0.4) / 0.8 * 90.0, 0.0, 90.0)
        return JointTargets(joints={"base": base, "shoulder": shoulder, "elbow": elbow})
