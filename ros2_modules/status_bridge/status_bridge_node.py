"""ROS2 status bridge replacing local status monitoring."""

from __future__ import annotations

import json
from typing import Any, Dict

import rclpy
from rclpy.node import Node
from std_msgs.msg import String

from ros2_modules.common.arm_protocol import utc_timestamp


class StatusBridgeNode(Node):
    """Publishes normalized robot status on ROS2 topic."""

    def __init__(self) -> None:
        super().__init__("status_bridge")
        self._latest_joints: Dict[str, float] = {"base": 0.0, "shoulder": 0.0, "elbow": 0.0}
        self._is_homed = False
        self._status_publisher = self.create_publisher(String, "/arm/status", 10)
        self.create_subscription(String, "/arm/joint_targets", self._on_joint_targets, 10)
        self._timer = self.create_timer(0.5, self._publish_status)
        self.get_logger().info("Status bridge started")

    def _on_joint_targets(self, msg: String) -> None:
        payload: Dict[str, Any] = json.loads(msg.data)
        joints = payload.get("joints", {})
        for name in self._latest_joints:
            if name in joints:
                self._latest_joints[name] = float(joints[name])
        if payload.get("mode") == "home":
            self._is_homed = True
        self._publish_status()

    def _publish_status(self) -> None:
        out = String()
        out.data = json.dumps(
            {
                "timestamp": utc_timestamp(),
                "is_homed": self._is_homed,
                "joints": self._latest_joints,
            },
            separators=(",", ":"),
        )
        self._status_publisher.publish(out)


def main() -> None:
    rclpy.init()
    node = StatusBridgeNode()
    try:
        rclpy.spin(node)
    finally:
        node.destroy_node()
        rclpy.shutdown()


if __name__ == "__main__":
    main()
