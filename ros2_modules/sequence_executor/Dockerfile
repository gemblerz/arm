FROM ros:humble-ros-base

WORKDIR /app
COPY ros2_modules /app/ros2_modules
ENV PYTHONPATH=/app

CMD ["/bin/bash", "-lc", "source /opt/ros/humble/setup.bash && python3 /app/ros2_modules/sequence_executor/sequence_executor_node.py"]
