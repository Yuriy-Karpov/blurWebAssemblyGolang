+ (function () {
    const img = new Image();
    img.src = './images/big-small-cat2.jpg';
    img.onload = function () {
        draw(this);
    };

    if (!WebAssembly.instantiateStreaming) { // polyfill
        WebAssembly.instantiateStreaming = async(resp, importObject) => {
            const source = await(await resp).arrayBuffer();
            return await WebAssembly.instantiate(source, importObject);
        };
    }

    const go = new Go();

    WebAssembly.instantiateStreaming(fetch("main.wasm"), go.importObject).then((result) => {
        go.run(result.instance);
    });

    function draw(img) {
        var canvas = document.getElementById('canvas');
        var canvas2 = document.getElementById("canvas2");
        var ctx = canvas.getContext('2d');
        var ctx2 = canvas2.getContext("2d");

        canvas.width = img.naturalWidth;
        canvas.height = img.naturalHeight;
        canvas2.width = img.naturalWidth;
        canvas2.height = img.naturalHeight;

        ctx.drawImage(img, 0, 0);
        ctx2.drawImage(img, 0, 0);
        // img.style.display = 'none';
        var imageData = ctx.getImageData(0, 0, canvas.width, canvas.height);
        var target = ctx.getImageData(0, 0, canvas.width, canvas.height);

        let rgba = imageData.data;

        const w = img.naturalWidth;
        const h = img.naturalHeight;
        const len = img.naturalWidth * img.naturalHeight;

        // let rgbaUnit = new Uint8Array(imageData.data.buffer);
        let rSrc = new Uint8Array(len);
        let gSrc = new Uint8Array(len);
        let bSrc = new Uint8Array(len);
        let aSrc = new Uint8Array(len);
        // let targetOut = new Uint8Array(len * 4);

        let offset = 0;

        for (let i = 0; i < len; i++) {
            rSrc[i] = rgba[offset++];
            gSrc[i] = rgba[offset++];
            bSrc[i] = rgba[offset++];
            aSrc[i] = rgba[offset++];
        }

        function runAsm() {
            console.time('runAsm')
            // This line call the webAssembly function
            const r = 10;
            // const promisR = new Promise((resolve, reject) => {
            //     calcAsm(rSrc, rSrc, w, h, r, (err, message) => {
            //         if (err) {
            //             console.error(err);
            //             return reject(err)
            //         }
            //         resolve(message);
            //     });
            // })
            // const promisG = new Promise((resolve, reject) => {
            //     calcAsm(gSrc, gSrc, w, h, r, (err, message) => {
            //         if (err) {
            //             console.error(err);
            //             return reject(err)
            //         }
            //         resolve(message);
            //     });
            // })
            // const promisB = new Promise((resolve, reject) => {
            //     calcAsm(bSrc, bSrc, w, h, r, (err, message) => {
            //         if (err) {
            //             console.error(err);
            //             return reject(err)
            //         }
            //         resolve(message);
            //     });
            // })
                calcAsm(rSrc, rSrc, w, h, r, (err, message) => {
                    if (err) {
                        console.error(err);
                    }
                    for (i = 0, offset = 0; i < len; i++) {
                        rgba[offset++] = rSrc[i];
                        rgba[offset++] = gSrc[i];
                        rgba[offset++] = bSrc[i];
                        rgba[offset++] = aSrc[i]; // or just increase offset if you skipped alpha
                    }
        
                    ctx2.putImageData(imageData, 0, 0);
                    console.timeEnd("runAsm")
                });
            // Promise.all([promisR, promisG, promisB]).then(value => {
            //     for (i = 0, offset = 0; i < len; i++) {
            //         rgba[offset++] = rSrc[i];
            //         rgba[offset++] = gSrc[i];
            //         rgba[offset++] = bSrc[i];
            //         rgba[offset++] = aSrc[i]; // or just increase offset if you skipped alpha
            //     }
    
            //     ctx2.putImageData(imageData, 0, 0);
            //     console.timeEnd("runAsm")
            // })

        }
        console.log('++END WebAssembly!');

        document
            .getElementById("control")
            .addEventListener("click", function (params) {
                runAsm()
            });
    }

})()
