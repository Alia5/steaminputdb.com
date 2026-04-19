import { build, defineConfig, type PluginOption } from 'vite';
import { svelte } from '@sveltejs/vite-plugin-svelte';

import { dirname, resolve } from 'node:path';
import { fileURLToPath } from 'node:url';
import { readdirSync, rmSync } from 'node:fs';
import svg from '@poppanator/sveltekit-svg';


const __dirname = dirname(fileURLToPath(import.meta.url));

const goTmpl = (): PluginOption => ({
    name: 'go-tmpl-param',
    generateBundle(_, bundle) {
        Object.values(bundle)
            .filter((c): c is typeof c & { type: 'chunk' } => c.type === 'chunk')
            .forEach((chunk) => {
                chunk.code = chunk.code.replace(/goTmpl\(["'](.+?)["']\)/g, '<<%$1%>>');
                chunk.code = chunk.code.replace(/,(<<%\s*(?:if\b|else\b|end\b)[^%]*%>>)/g, '$1');
                chunk.code = chunk.code.replace(/(<<%\s*(?:if\b|else\b|end\b)[^%]*%>>),/g, '$1');
            });
    }
});

let svgoPrefixIdsCount = 0;

if (!process.env.__VITE_CHILD_BUILD) {
    process.env.__VITE_CHILD_BUILD = '1';
    const inputs = readdirSync(resolve(__dirname, 'src/templates'))
        .filter((file) => file.endsWith('.ts'))
        .map((file) => [file.replace('.ts', ''), resolve(__dirname, 'src/templates', file)] as const);
    rmSync(resolve(__dirname, 'dist'), { recursive: true, force: true });
    await Promise.all(
        inputs.map(([name, path]) =>
            build({
                configFile: resolve(__dirname, 'vite.config.ts'),
                build: {
                    emptyOutDir: false,
                    rollupOptions: {
                        input: { [name]: path }
                    }
                }
            })
        )
    );
    process.exit(0);
}

export default defineConfig({
    resolve: {
        alias: {
            $lib: resolve(__dirname, 'src/lib')
        }
    },
    plugins: [
        svg({
            includePaths: [
                './src/lib/assets/'
            ],
            svgoOptions: {
                plugins: [
                    'preset-default',
                    {
                        name: 'prefixIds',
                        params: {
                            delim: '',
                            prefix: () => svgoPrefixIdsCount++
                        }
                    // eslint-disable-next-line @typescript-eslint/no-explicit-any
                    } as any
                ]
            }
        }),
        svelte(),
        goTmpl()
    ],
    build: {
        assetsInlineLimit: Infinity,
        rolldownOptions: {
            treeshake: true,
            output: {
                inlineDynamicImports: true,
                entryFileNames: '[name].js.tmpl',
                intro: '(function() {\nlet __INJECT_RETURN;\n',
                outro: '\n return __INJECT_RETURN;\n})();'
            }
        }
    }
});
