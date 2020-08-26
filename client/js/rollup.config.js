import typescript from 'rollup-plugin-typescript2';

export default {
    input: 'im.ts',
    plugins: [typescript()],
    output: [
        {
            file: 'dist/im.js',
            format: 'iife',
            name: 'EventStream',
        },
    ],
};
