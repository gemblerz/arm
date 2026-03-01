# ARM ROS2 Controller (Python)

This repository now provides a modular ROS2 (Python) control stack for the robot arm.

## What Changed

- Status monitoring is moved to ROS2 (`status-bridge` module).
- Kinematics control is handled by a ROS2 robot controller node (`kinematics-controller`).
- Sequence execution is revised for movement commands like:
  - `home`
  - `move_cartesian` (also accepts `move_cartician` for compatibility)
- Each module has its own Docker image (ROS2 Humble base) and can run independently in the same ROS2 network.

Legacy Go code is archived under `archieve/`, and the active control path for this revision is the ROS2 Python modules under `ros2_modules/`.

---

## Modules

### 1) Sequence Executor
Path: `ros2_modules/sequence_executor`

Responsibilities:
- Subscribes to `/arm/sequence_command`
- Routes commands to:
  - `/arm/home_command`
  - `/arm/cartesian_target`
- Publishes execution state to `/arm/sequence_status`

### 2) Kinematics Controller
Path: `ros2_modules/kinematics_controller`

Responsibilities:
- Subscribes to `/arm/cartesian_target` and `/arm/home_command`
- Computes joint targets using a ROS2 robot controller kinematics solver
- Publishes `/arm/joint_targets`

### 3) Status Bridge
Path: `ros2_modules/status_bridge`

Responsibilities:
- Subscribes to `/arm/joint_targets`
- Publishes normalized robot status on `/arm/status`

### 4) Display Controller
Path: `ros2_modules/display_controller`

Responsibilities:
- Subscribes to `/arm/status` and `/arm/sequence_status`
- Formats compact status lines for SSD1306-compatible display output
- Uses mock rendering by default, with optional `luma.oled` backend when available

---

## Quick Start (Docker)

### Prerequisites

- Docker
- Docker Compose plugin

### 1. Build and start all ROS2 modules

```bash
docker compose -f docker-compose.ros2.yml up --build
```

### 2. Optional: Set ROS domain

```bash
export ROS_DOMAIN_ID=42
docker compose -f docker-compose.ros2.yml up --build
```

---

## Send Commands

Open another terminal and publish commands into the ROS network:

### Home

```bash
docker compose -f docker-compose.ros2.yml exec sequence-executor \
  bash -lc "source /opt/ros/humble/setup.bash && \
  ros2 topic pub --once /arm/sequence_command std_msgs/msg/String '{data: \"{\\\"command\\\":\\\"home\\\"}\"}'"
```

### Move to Cartesian coordinate

```bash
docker compose -f docker-compose.ros2.yml exec sequence-executor \
  bash -lc "source /opt/ros/humble/setup.bash && \
  ros2 topic pub --once /arm/sequence_command std_msgs/msg/String '{data: \"{\\\"command\\\":\\\"move_cartesian\\\",\\\"target\\\":{\\\"x\\\":0.25,\\\"y\\\":0.10,\\\"z\\\":0.30}}\"}'"
```

---

## Observe Topics

```bash
docker compose -f docker-compose.ros2.yml exec status-bridge \
  bash -lc "source /opt/ros/humble/setup.bash && ros2 topic echo /arm/status"
```

Other useful topics:

- `/arm/sequence_status`
- `/arm/cartesian_target`
- `/arm/home_command`
- `/arm/joint_targets`

---

## Local Tests (Python modules)

Run focused unit tests for the revised ROS2 modules:

```bash
PYTHONPATH=. python3 -m unittest discover -s ros2_modules/tests -v
```

---

## File Layout

```text
ros2_modules/
  common/
  sequence_executor/
  kinematics_controller/
  status_bridge/
  display_controller/
  tests/
docker-compose.ros2.yml
```
