import * as THREE from 'three';
import { OrbitControls } from 'https://unpkg.com/three@0.160.0/examples/jsm/controls/OrbitControls.js';

let scene, camera, renderer, controls;
const joints = {}; // Store references to the arm's joints

function init() {
    // Scene
    scene = new THREE.Scene();
    scene.background = new THREE.Color(0x222222);

    // Camera
    camera = new THREE.PerspectiveCamera(75, window.innerWidth / window.innerHeight, 0.1, 1000);
    camera.position.set(25, 35, 50);

    // Renderer
    renderer = new THREE.WebGLRenderer({ antialias: true });
    renderer.setSize(window.innerWidth, window.innerHeight);
    document.body.appendChild(renderer.domElement);

    // Lights
    const ambientLight = new THREE.AmbientLight(0xffffff, 1.5);
    scene.add(ambientLight);

    const directionalLight = new THREE.DirectionalLight(0xffffff, 1.5);
    directionalLight.position.set(20, 30, 40);
    scene.add(directionalLight);

    // Controls
    controls = new OrbitControls(camera, renderer.domElement);
    controls.enableDamping = true;

    // Grid Helper
    const gridHelper = new THREE.GridHelper(100, 20);
    scene.add(gridHelper);

    // Robot Arm
    createRobotArm();

    // WebSocket Connection (Placeholder)
    connectWebSocket();

    // Animation Loop
    animate();

    // Handle window resize
    window.addEventListener('resize', onWindowResize, false);
}

function createRobotArm() {
    const armMaterial = new THREE.MeshStandardMaterial({ color: 0x888888 });
    const jointMaterial = new THREE.MeshStandardMaterial({ color: 0xcccccc });

    // Base
    const base = new THREE.Mesh(new THREE.CylinderGeometry(5, 5, 2, 32), armMaterial);
    scene.add(base);

    // Joint 1 (Pan)
    joints.pan = new THREE.Object3D();
    joints.pan.position.y = 2;
    base.add(joints.pan);

    const link1 = new THREE.Mesh(new THREE.BoxGeometry(4, 8, 4), armMaterial);
    link1.position.y = 4;
    joints.pan.add(link1);

    // Joint 2 (Tilt)
    joints.tilt = new THREE.Object3D();
    joints.tilt.position.y = 8;
    joints.pan.add(joints.tilt);

    const link2 = new THREE.Mesh(new THREE.BoxGeometry(15, 3, 3), armMaterial);
    link2.position.x = 7.5;
    joints.tilt.add(link2);

    // Joint 3
    joints.joint3 = new THREE.Object3D();
    joints.joint3.position.x = 15;
    joints.tilt.add(joints.joint3);

    const link3 = new THREE.Mesh(new THREE.BoxGeometry(12, 2.5, 2.5), armMaterial);
    link3.position.x = 6;
    joints.joint3.add(link3);

    // Joint 4
    joints.joint4 = new THREE.Object3D();
    joints.joint4.position.x = 12;
    joints.joint3.add(joints.joint4);

    const link4 = new THREE.Mesh(new THREE.BoxGeometry(3, 2, 2), armMaterial);
    link4.position.x = 1.5;
    joints.joint4.add(link4);

    // Gripper (simplified)
    const gripperBase = new THREE.Object3D();
    gripperBase.position.x = 3;
    joints.joint4.add(gripperBase);

    const gripperLeft = new THREE.Mesh(new THREE.BoxGeometry(0.5, 1, 3), jointMaterial);
    gripperLeft.position.z = -1.5;
    gripperBase.add(gripperLeft);

    const gripperRight = new THREE.Mesh(new THREE.BoxGeometry(0.5, 1, 3), jointMaterial);
    gripperRight.position.z = 1.5;
    gripperBase.add(gripperRight);

    console.log("Robot arm model created.");
}

function updateJoints(angles) {
    if (!angles) return;

    // Note: The conversion from motor steps to radians will happen here.
    // For now, we'll just use the incoming values as if they are radians.
    if (joints.pan) joints.pan.rotation.y = angles.pan || 0;
    if (joints.tilt) joints.tilt.rotation.z = angles.tilt || 0;
    if (joints.joint3) joints.joint3.rotation.z = angles.joint3 || 0;
    if (joints.joint4) joints.joint4.rotation.y = angles.joint4 || 0;
}

function connectWebSocket() {
    const socket = new WebSocket("ws://localhost:8080/ws");

    socket.onopen = function(event) {
        console.log("WebSocket connection established.");
        socket.send("Hello Server!");
    };

    socket.onmessage = function(event) {
        try {
            const data = JSON.parse(event.data);
            if (data.jointPositions) {
                // Convert steps to radians for visualization
                const angles = {
                    pan: (data.jointPositions.pan || 0) * (Math.PI / 1800), // Example conversion
                    tilt: (data.jointPositions.tilt || 0) * (Math.PI / 1800),
                    joint3: (data.jointPositions.joint3 || 0) * (Math.PI / 1800),
                    joint4: (data.jointPositions.joint4 || 0) * (Math.PI / 1800)
                };
                updateJoints(angles);
            }
        } catch (e) {
            console.error("Error parsing WebSocket message:", e);
        }
    };

    socket.onclose = function(event) {
        console.log("WebSocket connection closed. Attempting to reconnect...");
        setTimeout(connectWebSocket, 2000); // Reconnect after 2 seconds
    };

    socket.onerror = function(error) {
        console.error("WebSocket error:", error);
    };
}

function onWindowResize() {
    camera.aspect = window.innerWidth / window.innerHeight;
    camera.updateProjectionMatrix();
    renderer.setSize(window.innerWidth, window.innerHeight);
}

function animate() {
    requestAnimationFrame(animate);
    controls.update();
    renderer.render(scene, camera);
}

init();
