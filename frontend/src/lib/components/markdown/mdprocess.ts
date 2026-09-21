import hljs from 'highlight.js';
import rehypeExternalLinks from 'rehype-external-links';
import rehypeHighlight from 'rehype-highlight';
import rehypeRaw from 'rehype-raw';
import rehypeSanitize, { defaultSchema } from 'rehype-sanitize';
import rehypeStringify from 'rehype-stringify';
import remarkGfm from 'remark-gfm';
import remarkParse from 'remark-parse';
import remarkRehype from 'remark-rehype';
import { unified, type PluggableList } from 'unified';
import { rehypeWrap } from './rehypeWrap';


export const processMarkdown = (extraPlugins: PluggableList = []) => {

    const result = unified()
        .use(remarkParse, { fragment: true })
        .use(remarkGfm)
        .use(remarkRehype, { allowDangerousHtml: true })
        .use(rehypeRaw)
        .use(extraPlugins)
        .use(rehypeSanitize, {
            ...defaultSchema,
            tagNames: [...(defaultSchema.tagNames ?? []), 'citation-ref'],
            attributes: {
                ...defaultSchema.attributes,
                'citation-ref': ['idx', 'filepath', 'page', 'content', 'chatId'],
                'div': [...(defaultSchema.attributes?.div ?? []), ['className', 'math', 'math-display']],
                'span': [...(defaultSchema.attributes?.span ?? []), ['className', 'math', 'math-inline']],
                'code': [...(defaultSchema.attributes?.code ?? []), ['className', 'math-inline', 'math-display']]
            }
        })
        .use(rehypeExternalLinks, { target: '_blank', rel: ['noopener'] })
        .use(rehypeWrap, ['table'])
        .use(rehypeHighlight, { detect: true })
        .use(rehypeStringify);
    return result;
};


export const initCodeHighlight = async () => {
    hljs.configure({
        ignoreUnescapedHTML: true
    });
};
