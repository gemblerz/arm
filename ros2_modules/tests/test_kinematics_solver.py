import unittest

from ros2_modules.common.arm_protocol import CartesianTarget
from ros2_modules.kinematics_controller.kinematics_solver import RobotControllerKinematics


class TestRobotControllerKinematics(unittest.TestCase):
    def test_solve_nominal_target(self) -> None:
        solver = RobotControllerKinematics(max_reach_xy=1.0)
        result = solver.solve(CartesianTarget(x=0.5, y=0.5, z=0.0))
        self.assertAlmostEqual(result.joints["base"], 45.0, places=1)
        self.assertGreaterEqual(result.joints["shoulder"], 0.0)
        self.assertLessEqual(result.joints["shoulder"], 90.0)
        self.assertEqual(result.joints["elbow"], 45.0)

    def test_solve_clamps_radius(self) -> None:
        solver = RobotControllerKinematics(max_reach_xy=0.8)
        result = solver.solve(CartesianTarget(x=2.0, y=0.0, z=0.8))
        self.assertEqual(result.joints["base"], 0.0)
        self.assertEqual(result.joints["shoulder"], 90.0)
        self.assertEqual(result.joints["elbow"], 90.0)

    def test_solve_clamps_low_z(self) -> None:
        solver = RobotControllerKinematics(max_reach_xy=0.8)
        result = solver.solve(CartesianTarget(x=0.1, y=0.1, z=-1.0))
        self.assertEqual(result.joints["elbow"], 0.0)


if __name__ == "__main__":
    unittest.main()
