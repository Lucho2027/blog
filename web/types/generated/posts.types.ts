/*
 * AUTO-GENERATED (types-only) from sqlc output.
 * Do not edit by hand.
 */

export interface CreateDraftPostArgs {
    title: string;
    excerpt: string | null;
    contentMd: string;
    userId: string;
}

export interface CreateDraftPostRow {
    id: string;
    title: string;
    excerpt: string | null;
    contentMd: string;
    status: string;
    userId: string;
    publishedAt: Date | null;
    createdAt: Date | null;
    updatedAt: Date | null;
}

export interface GetPostByIdArgs {
    id: string;
}

export interface GetPostByIdRow {
    id: string;
    title: string;
    excerpt: string | null;
    contentMd: string;
    status: string;
    userId: string;
    publishedAt: Date | null;
    createdAt: Date | null;
    updatedAt: Date | null;
}

export interface GetAllPostsArgs {
    offsetVal: string;
    limitVal: string;
}

export interface GetAllPostsRow {
    id: string;
    title: string;
    excerpt: string | null;
    contentMd: string;
    status: string;
    userId: string;
    publishedAt: Date | null;
    createdAt: Date | null;
    updatedAt: Date | null;
}

