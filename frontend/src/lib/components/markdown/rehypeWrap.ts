import type { Element, Root } from 'hast';
import type { Plugin } from 'unified';

type ParentNode = Root | Element;

const normalizeSelectors = (selectors: string[] = ['table']): string[] => selectors
    .map((value) => value.trim().toLowerCase())
    .filter(Boolean);

const wrapMatches = (node: ParentNode, selectors: string[], wrapperTagName: string): void => {
    node.children.forEach((child, index) => {
        if (child.type !== 'element') {
            return;
        }

        if (selectors.includes(child.tagName.toLowerCase())) {
            node.children[index] = {
                type: 'element',
                tagName: wrapperTagName,
                properties: {},
                children: [child]
            };
            return;
        }

        wrapMatches(child, selectors, wrapperTagName);
    });
};

export const rehypeWrap: Plugin<[(string[] | undefined)?, (string | undefined)?], Root> = (
    selectors: string[] = [],
    wrapperTagName = 'div'
) => {
    const normalizedSelectors = normalizeSelectors(selectors);

    return (tree) => {
        wrapMatches(tree, normalizedSelectors, wrapperTagName);
    };
};
