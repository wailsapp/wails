// The document model: a tree of nodes. Every node has a component type and a
// bag of props; container components also carry children. The tree is plain
// JSON so it round-trips through the Go LayoutService unchanged.

export type NodeId = string;

export interface UINode {
    id: NodeId;
    type: string;
    props: Record<string, unknown>;
    children?: UINode[];
}

export type ArtboardTheme = 'light' | 'dark';

export interface UIDocument {
    version: 1;
    name: string;
    theme: ArtboardTheme;
    root: UINode;
}

export const ROOT_ID = 'root';

let counter = 0;

/** Short, unique-enough ids: a time slice plus a per-session counter. */
export function newId(): NodeId {
    counter += 1;
    return (Date.now().toString(36).slice(-4) + counter.toString(36)).toLowerCase();
}

export function emptyDocument(name = 'Untitled'): UIDocument {
    return {
        version: 1,
        name,
        theme: 'light',
        root: {id: ROOT_ID, type: 'root', props: {}, children: []},
    };
}

/** Depth-first search for a node by id. */
export function findNode(root: UINode, id: NodeId): UINode | null {
    if (root.id === id) {
        return root;
    }
    for (const child of root.children ?? []) {
        const hit = findNode(child, id);
        if (hit) {
            return hit;
        }
    }
    return null;
}

/** Returns the parent of a node and the index of the node within it. */
export function findParent(root: UINode, id: NodeId): { parent: UINode; index: number } | null {
    for (const [index, child] of (root.children ?? []).entries()) {
        if (child.id === id) {
            return {parent: root, index};
        }
        const hit = findParent(child, id);
        if (hit) {
            return hit;
        }
    }
    return null;
}

/** Ancestors from the root down to (but excluding) the node. */
export function pathTo(root: UINode, id: NodeId): UINode[] {
    const path: UINode[] = [];
    const walk = (node: UINode): boolean => {
        if (node.id === id) {
            return true;
        }
        for (const child of node.children ?? []) {
            path.push(node);
            if (walk(child)) {
                return true;
            }
            path.pop();
        }
        return false;
    };
    return walk(root) ? path : [];
}

/** True when `maybeAncestorId` is `id` itself or one of its ancestors. */
export function isWithin(root: UINode, id: NodeId, maybeAncestorId: NodeId): boolean {
    if (id === maybeAncestorId) {
        return true;
    }
    return pathTo(root, id).some((n) => n.id === maybeAncestorId);
}

/** Deep clone a node with fresh ids (for duplicate / paste). */
export function cloneNode(node: UINode): UINode {
    return {
        id: newId(),
        type: node.type,
        props: {...node.props},
        children: node.children?.map(cloneNode),
    };
}

export function walk(node: UINode, fn: (node: UINode, depth: number) => void, depth = 0): void {
    fn(node, depth);
    for (const child of node.children ?? []) {
        walk(child, fn, depth + 1);
    }
}

export function countNodes(node: UINode): number {
    let n = 0;
    walk(node, () => { n += 1; });
    return n - 1; // exclude the root
}
