export enum UserRole {
    GUIDE = 'guide',
    TOURIST = 'tourist', 
    ADMINISTRATOR = 'administrator'
}

export interface RoleOption {
    value: UserRole;
    label: string;
}

export const ROLE_OPTIONS: RoleOption[] = [
    { value: UserRole.GUIDE, label: 'Vodič' },
    { value: UserRole.TOURIST, label: 'Turista' }
];