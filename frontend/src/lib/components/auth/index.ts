import type { Role } from '$lib/utils/auth-guards';

export { default as RoleGate } from './RoleGate.svelte';

export const ROLES: Record<string, Role> = {
  owner: 'OWNER',
  admin: 'ADMIN',
  manager: 'MANAGER',
  viewer: 'VIEWER',
} as const;
