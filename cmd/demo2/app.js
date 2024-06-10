'use strict';

var RAY_COUNT = 1;

var container;
var camera, scene, renderer, controls;

var modelObject;
var intersectingTriangles;
var rayLines;

var windowHalfX = window.innerWidth / 2;
var windowHalfY = window.innerHeight / 2;

var bvh; // keeps the bvh data structure, build from the model

var hexColors = [
    0xFF1493,
    0xB23AEE,
    0x1E90FF,
    0x00F5FF,
    0x00FF7F,
    0x7FFF00,
    0xCDCD00,
    0xFFA500
];

var randomValues = [];

var raycaster = new THREE.Raycaster();
var mouse = new THREE.Vector2(-10000, -10000);

function onMouseMove(event) {

// calculate mouse position in normalized device coordinates
// (-1 to +1) for both components

    mouse.x = (event.clientX / window.innerWidth) * 2 - 1;
    mouse.y = -(event.clientY / window.innerHeight) * 2 + 1;

}

function init() {

    container = document.createElement('div');
    document.body.appendChild(container);

    window.addEventListener('mousemove', onMouseMove, false);

    camera = new THREE.PerspectiveCamera(45, window.innerWidth / window.innerHeight, 1, 2000);
    camera.position.z = -325;
    camera.position.y = -75;

    controls = new THREE.OrbitControls(camera);

// scene
    scene = new THREE.Scene();

    var ambient = new THREE.AmbientLight(0x101030);
    scene.add(ambient);

    var directionalLight0 = new THREE.DirectionalLight(0xffeedd);
    directionalLight0.position.set(0, 0, 1);
    scene.add(directionalLight0);

    var directionalLight1 = new THREE.DirectionalLight(0xffeedd);
    directionalLight1.position.set(0, 0, -1);
    scene.add(directionalLight1);

    for (var i = 0; i < 5000; i++) {
        randomValues.push(Math.random());
    }

    var manager = new THREE.LoadingManager();
    manager.onProgress = function (item, loaded, total) {

        console.log(item, loaded, total);

    };

    var onProgress = function (xhr) {
        if (xhr.lengthComputable) {
            var percentComplete = xhr.loaded / xhr.total * 100;
            console.log(Math.round(percentComplete, 2) + '% downloaded');
        }
    };

    var onError = function (xhr) {
    };

    var loader = new THREE.OBJLoader(manager);
    loader.load('Armadillo_100k.obj', function (object) {
        console.log("finished loading model from file");
        object.traverse(function (child) {

            if (child instanceof THREE.Mesh) {

                buildBVH(child);
                child.geometry.computeVertexNormals();
            }

        });

        scene.add(object);
        modelObject = object;

    }, onProgress, onError);

//

    renderer = new THREE.WebGLRenderer();
    renderer.setPixelRatio(window.devicePixelRatio);
    renderer.setSize(window.innerWidth, window.innerHeight);
    container.appendChild(renderer.domElement);

    document.addEventListener('mousemove', onDocumentMouseMove, false);
    window.addEventListener('resize', onWindowResize, false);
}

function onWindowResize() {

    windowHalfX = window.innerWidth / 2;
    windowHalfY = window.innerHeight / 2;

    camera.aspect = window.innerWidth / window.innerHeight;
    camera.updateProjectionMatrix();

    renderer.setSize(window.innerWidth, window.innerHeight);

}

function onDocumentMouseMove(event) {

}

function animate() {

    requestAnimationFrame(animate);
    update();
    render();

}

function update() {
    if (intersectingTriangles) {
        scene.remove(intersectingTriangles);
        scene.remove(rayLines);
    }

    intersectingTriangles = new THREE.Object3D();
    rayLines = new THREE.Object3D();

    if (!bvh) {
        return;
    }

    for (var i = 0; i < RAY_COUNT; i++) {
        var rayOrigin = mouse.clone();
        raycaster.setFromCamera(rayOrigin, camera);

        // raycaster.ray.origin.x += (randomValues[(i * 6) % randomValues.length] - 0.5) * 10;
        // raycaster.ray.origin.y += (randomValues[(i * 6 + 1) % randomValues.length] - 0.5) * 10;
        // raycaster.ray.origin.z += (randomValues[(i * 6 + 2) % randomValues.length] - 0.5) * 10;
        // raycaster.ray.direction.x += (randomValues[(i * 6 + 3) % randomValues.length] - 0.5) * 0.01;
        // raycaster.ray.direction.y += (randomValues[(i * 6 + 4) % randomValues.length] - 0.5) * 0.01;
        // raycaster.ray.direction.z += (randomValues[(i * 6 + 5) % randomValues.length] - 0.5) * 0.01;
        raycaster.ray.direction.normalize();

        var rayColor = hexColors[i % hexColors.length];

        // console.log("intersectRay", raycaster.ray.origin, raycaster.ray.direction, true);
        var intersectResult = bvh.intersectRay(raycaster.ray.origin, raycaster.ray.direction, true);
        if (intersectResult.length > 0) {
            intersectResult.forEach(function (intersection) {
                var triangle = intersection.triangle;
                var geometry = new THREE.Geometry();

                geometry.vertices.push(triangle[0]);
                geometry.vertices.push(triangle[1]);
                geometry.vertices.push(triangle[2]);

                geometry.faces.push(new THREE.Face3(0, 1, 2));

                var redMat = new THREE.MeshBasicMaterial({color: rayColor, side: THREE.DoubleSide});
                var intersectingTriangle = new THREE.Mesh(geometry, redMat);
                intersectingTriangles.add(intersectingTriangle);
            });
        }
    }

    intersectingTriangles.frustumCulled = false;
    scene.add(intersectingTriangles);
}

function render() {
    renderer.render(scene, camera);

}

function buildBVH(threeMesh) {
    var vertexArr = threeMesh.geometry.attributes.position.array;
    // var triangleCount = vertexArr.length / 9;
    // var triangles = [];
    // triangles.length = triangleCount;
    //
    // for (var i = 0; i < triangleCount; i++) {
    //     triangles[i] = [
    //         {x: vertexArr[i * 9], y: vertexArr[i * 9 + 1], z: vertexArr[i * 9 + 2]},
    //         {x: vertexArr[i * 9 + 3], y: vertexArr[i * 9 + 4], z: vertexArr[i * 9 + 5]},
    //         {x: vertexArr[i * 9 + 6], y: vertexArr[i * 9 + 7], z: vertexArr[i * 9 + 8]}
    //     ]
    // }

    var start = new Date().getTime();
    // bvh = new bvhtree.BVH(triangles, 7);
    bvh = bvhtree.bvhFromVertexArray(vertexArr, 7);
    console.log("input triangles: ", threeMesh.geometry.attributes.position.array.length / 9);
    console.log("building bvh took ", new Date().getTime() - start, " ms");


// hide loader animation
    document.getElementById('loading-wrapper').style.display = 'none';
}