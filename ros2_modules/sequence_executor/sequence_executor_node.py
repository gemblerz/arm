"""ROS2 sequence executor for home and Cartesian commands."""

from __future__ import annotations

import json

import rclpy
from rclpy.node import Node
from std_msgs.msg import String

from ros2_modules.common.arm_protocol import utc_timestamp
from ros2_modules.sequence_executor.sequence_parser import parse_command


class SequenceExecutorNode(Node):
    """Consumes sequence commands and routes ROS2 control topics."""

    def __init__(self) -> None:
        super().__init__("sequence_executor")
        self._home_pub = self.create_publisher(String, "/arm/home_command", 10)
        self._target_pub = self.create_publisher(String, "/arm/cartesian_target", 10)
        self._status_pub = self.create_publisher(String, "/arm/sequence_status", 10)
        self.create_subscription(String, "/arm/sequence_command", self._on_command, 10)
        self.get_logger().info("Sequence executor started")

    def _on_command(self, msg: String) -> None:
        payload = json.loads(msg.data)
        parsed = parse_command(payload)
        if parsed.action == "home":
            command = String()
            command.data = json.dumps({"home": True}, separators=(",", ":"))
            self._home_pub.publish(command)
            self._publish_status("home", "accepted")
            return

        if parsed.target is not None:
            target = String()
            target.data = json.dumps(
                {"x": parsed.target.x, "y": parsed.target.y, "z": parsed.target.z},
                separators=(",", ":"),
            )
            self._target_pub.publish(target)
            self._publish_status("move_cartesian", "accepted")

    def _publish_status(self, action: str, state: str) -> None:
        status = String()
        status.data = json.dumps(
            {"timestamp": utc_timestamp(), "action": action, "state": state},
            separators=(",", ":"),
        )
        self._status_pub.publish(status)


def main() -> None:
    rclpy.init()
    node = SequenceExecutorNode()
    try:
        rclpy.spin(node)
    finally:
        node.destroy_node()
        rclpy.shutdown()


if __name__ == "__main__":
    main()
