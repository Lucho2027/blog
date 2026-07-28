import { Sql } from "postgres";

export const createRefreshTokenQuery = `-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (token, created_at, updated_at, user_id, expires_at, revoked_at)
VALUES ($1, now(), now(), $2, $3, null)
RETURNING token, created_at, updated_at, user_id, expires_at, revoked_at`;

export interface CreateRefreshTokenArgs {
    token: string;
    userId: string;
    expiresAt: Date | null;
}

export interface CreateRefreshTokenRow {
    token: string;
    createdAt: Date;
    updatedAt: Date;
    userId: string;
    expiresAt: Date | null;
    revokedAt: Date | null;
}

export async function createRefreshToken(sql: Sql, args: CreateRefreshTokenArgs): Promise<CreateRefreshTokenRow | null> {
    const rows = await sql.unsafe(createRefreshTokenQuery, [args.token, args.userId, args.expiresAt]).values();
    if (rows.length !== 1) {
        return null;
    }
    const row = rows[0];
    return {
        token: row[0],
        createdAt: row[1],
        updatedAt: row[2],
        userId: row[3],
        expiresAt: row[4],
        revokedAt: row[5]
    };
}

export const getUserIDFromRefreshTokenQuery = `-- name: GetUserIDFromRefreshToken :one
Select user_id from refresh_tokens
where token = $1
and revoked_at IS NULL and expires_at > now()`;

export interface GetUserIDFromRefreshTokenArgs {
    token: string;
}

export interface GetUserIDFromRefreshTokenRow {
    userId: string;
}

export async function getUserIDFromRefreshToken(sql: Sql, args: GetUserIDFromRefreshTokenArgs): Promise<GetUserIDFromRefreshTokenRow | null> {
    const rows = await sql.unsafe(getUserIDFromRefreshTokenQuery, [args.token]).values();
    if (rows.length !== 1) {
        return null;
    }
    const row = rows[0];
    return {
        userId: row[0]
    };
}

export const revokeRefreshTokenQuery = `-- name: RevokeRefreshToken :exec
update refresh_tokens set updated_at = now(), revoked_at = now()
where token = $1`;

export interface RevokeRefreshTokenArgs {
    token: string;
}

export async function revokeRefreshToken(sql: Sql, args: RevokeRefreshTokenArgs): Promise<void> {
    await sql.unsafe(revokeRefreshTokenQuery, [args.token]);
}

