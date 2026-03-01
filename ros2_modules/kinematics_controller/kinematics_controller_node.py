"""ROS2 node that maps Cartesian targets to joint targets."""

from __future__ import annotations

import json

import rclpy
from rclpy.node import Node
from std_msgs.msg import String

from ros2_modules.common.arm_protocol import CartesianTarget, utc_timestamp
from ros2_modules.kinematics_controller.kinematics_solver import RobotControllerKinematics


class KinematicsControllerNode(Node):
    """Robot controller node providing kinematics based joint targets."""

    def __init__(self) -> None:
        super().__init__("kinematics_controller")
        self._solver = RobotControllerKinematics()
        self._publisher = self.create_publisher(String, "/arm/joint_targets", 10)
        self._subscriber = self.create_subscription(
            String, "/arm/cartesian_target", self._on_target, 10
        )
        self._home_subscriber = self.create_subscription(
            String, "/arm/home_command", self._on_home, 10
        )
        self.get_logger().info("Kinematics controller started")

    def _on_target(self, msg: String) -> None:
        payload = json.loads(msg.data)
        target = CartesianTarget(
            x=float(payload["x"]),
            y=float(payload["y"]),
            z=float(payload["z"]),
        )
        joints = self._solver.solve(target)
        self._publish_joint_targets(joints.joints, mode="cartesian")

    def _on_home(self, _msg: String) -> None:
        self._publish_joint_targets({"base": 0.0, "shoulder": 0.0, "elbow": 0.0}, mode="home")

    def _publish_joint_targets(self, joints: dict[str, float], mode: str) -> None:
        out = String()
        out.data = json.dumps(
            {"timestamp": utc_timestamp(), "mode": mode, "joints": joints},
            separators=(",", ":"),
        )
        self._publisher.publish(out)


def main() -> None:
    rclpy.init()
    node = KinematicsControllerNode()
    try:
        rclpy.spin(node)
    finally:
        node.destroy_node()
        rclpy.shutdown()


if __name__ == "__main__":
    main()
