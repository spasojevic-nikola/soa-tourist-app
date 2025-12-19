import { UserRole } from './role.enum';

export interface Registration {
    name: string,
    surname: string,
    email: string,
    username: string,
    password: string,
    role: UserRole
}