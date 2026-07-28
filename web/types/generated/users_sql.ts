import { Sql } from "postgres";

export const createUserQuery = `-- name: CreateUser :one
insert into users(email, hash_pw)
values ($1, $2)
returning id`;

export interface CreateUserArgs {
    email: string;
    hashPw: string;
}

export interface CreateUserRow {
    id: string;
}

export async function createUser(sql: Sql, args: CreateUserArgs): Promise<CreateUserRow | null> {
    const rows = await sql.unsafe(createUserQuery, [args.email, args.hashPw]).values();
    if (rows.length !== 1) {
        return null;
    }
    const row = rows[0];
    return {
        id: row[0]
    };
}

