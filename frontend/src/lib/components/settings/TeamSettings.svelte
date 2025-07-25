<script lang="ts">
  import type { ModelsTeamMember } from "$lib/api/api";
  import { RoleGate, ROLES } from "$lib/components/auth";
  import { Button, Input, Modal } from "$lib/components/ui";
  import { auth } from "$lib/stores/auth";
  import type { Role } from "$lib/utils/auth-guards";
  import { formatDate } from "$lib/utils/date-utils";
  import {
    Eye,
    Loader2,
    Mail,
    Settings,
    Shield,
    Trash2,
    UserPlus,
    Users,
  } from "lucide-svelte";

  interface Props {
    onSuccess?: (message: string) => void;
    onError?: (message: string) => void;
  }

  let { onSuccess, onError }: Props = $props();

  // Component state
  let loading = $state(false);
  let teamMembers = $state<ModelsTeamMember[]>([]);
  let showInviteModal = $state(false);
  let inviteEmail = $state("");
  let inviteRole = $state<string>("VIEWER");
  let inviting = $state(false);
  let removingMemberId = $state<string | null>(null);
  let resendingInvitationId = $state<string | null>(null);

  // Load team members when component mounts and user is authenticated
  $effect(() => {
    if ($auth.isAuthenticated) {
      loadTeamMembers();
    }
  });

  async function loadTeamMembers() {
    loading = true;
    try {
      // Get API instance from auth store
      const api = auth.getApi();

      const response = await api.api.v1TeamMembersList();

      if (response.data.success && response.data.data) {
        teamMembers = response.data.data as ModelsTeamMember[];
      } else {
        throw new Error("Failed to load team members");
      }
    } catch (error) {
      console.error("Error loading team members:", error);
      onError?.("Failed to load team members");
    } finally {
      loading = false;
    }
  }

  async function inviteMember() {
    if (!inviteEmail || inviting) return;

    inviting = true;
    try {
      // Get API instance from auth store
      const api = auth.getApi();

      const response = await api.api.v1TeamMembersInviteCreate({
        email: inviteEmail,
        role: inviteRole,
      } as any);

      if (response.data.success) {
        onSuccess?.(`Invitation sent to ${inviteEmail}`);
        showInviteModal = false;
        inviteEmail = "";
        inviteRole = "VIEWER";

        // Reload team members
        await loadTeamMembers();
      } else {
        throw new Error("Failed to send invitation");
      }
    } catch (error: any) {
      const errorMessage =
        error.response?.data?.error?.message || "Failed to send invitation";
      onError?.(errorMessage);
    } finally {
      inviting = false;
    }
  }

  async function resendInvitation(invitationId: string) {
    if (resendingInvitationId) return;

    resendingInvitationId = invitationId;
    try {
      // Get API instance from auth store
      const api = auth.getApi();

      const response =
        await api.api.v1TeamMembersResendInvitationCreate(invitationId);

      if (response.data.success) {
        onSuccess?.("Invitation resent successfully");
      } else {
        throw new Error("Failed to resend invitation");
      }
    } catch (error: any) {
      const errorMessage =
        error.response?.data?.error?.message || "Failed to resend invitation";
      onError?.(errorMessage);
    } finally {
      resendingInvitationId = null;
    }
  }

  async function removeMember(memberId: string) {
    if (removingMemberId) return;

    removingMemberId = memberId;
    try {
      // Get API instance from auth store
      const api = auth.getApi();

      const response = await api.api.v1TeamMembersDelete(memberId);

      if (response.data.success) {
        onSuccess?.("Team member removed successfully");

        // Reload team members
        await loadTeamMembers();
      } else {
        throw new Error("Failed to remove team member");
      }
    } catch (error: any) {
      const errorMessage =
        error.response?.data?.error?.message || "Failed to remove team member";
      onError?.(errorMessage);
    } finally {
      removingMemberId = null;
    }
  }

  async function updateRole(memberId: string, newRole: Role) {
    try {
      // Get API instance from auth store
      const api = auth.getApi();

      const response = await api.api.v1TeamMembersRoleUpdate(memberId, {
        role: newRole,
      } as any);

      if (response.data.success) {
        onSuccess?.("Role updated successfully");

        // Reload team members
        await loadTeamMembers();
      } else {
        throw new Error("Failed to update role");
      }
    } catch (error: any) {
      const errorMessage =
        error.response?.data?.error?.message || "Failed to update role";
      onError?.(errorMessage);
    }
  }

  // Get role badge color
  function getRoleBadgeClass(role: string) {
    switch (role) {
      case ROLES.owner:
        return "bg-purple-100 text-purple-800";
      case ROLES.admin:
        return "bg-blue-100 text-blue-800";
      case ROLES.manager:
        return "bg-green-100 text-green-800";
      case ROLES.viewer:
        return "bg-gray-100 text-gray-800";
      default:
        return "bg-gray-100 text-gray-800";
    }
  }

  // Get role icon
  function getRoleIcon(role: string) {
    switch (role) {
      case ROLES.owner:
        return Shield;
      case ROLES.admin:
        return Settings;
      case ROLES.manager:
        return Users;
      case ROLES.viewer:
        return Eye;
      default:
        return Users;
    }
  }

  // Get user's role
  let currentUserRole = $derived(
    teamMembers.find((m: any) => {
      // For owner, check if account_id matches since member data might be invalid
      if (m.role === ROLES.owner && m.account_id === $auth.user?.id) {
        return true;
      }
      // For other roles, check email match
      return m.member?.email === $auth.user?.email;
    })?.role || ROLES.viewer,
  );

  // Check if user can manage team
  let canManageTeam = $derived(
    currentUserRole === ROLES.owner || currentUserRole === ROLES.admin,
  );
</script>

<div>
  <div class="mb-8">
    <div class="flex items-center justify-between">
      <div>
        <h2 class="text-2xl font-bold text-gray-900">Team Members</h2>
        <p class="mt-1 text-sm text-gray-600">
          Manage your team and their access levels
        </p>
      </div>
      <RoleGate roles={[ROLES.owner, ROLES.admin]}>
        <Button
          variant="gradient"
          size="sm"
          onclick={() => (showInviteModal = true)}
        >
          <UserPlus class="h-4 w-4 mr-2" />
          Invite Member
        </Button>
      </RoleGate>
    </div>
  </div>

  {#if loading}
    <div class="flex items-center justify-center py-12">
      <Loader2 class="h-8 w-8 animate-spin text-gray-400" />
    </div>
  {:else}
    <div class="space-y-4">
      {#each teamMembers as member}
        {@const RoleIcon = getRoleIcon(member.role || "")}
        <div class="bg-white border border-gray-200 rounded-lg p-4">
          <div class="flex items-center justify-between">
            <div class="flex items-center space-x-4">
              <div
                class="h-10 w-10 bg-gradient-to-br from-blue-500 to-purple-600 rounded-full flex items-center justify-center text-white font-semibold"
              >
                {member.member?.first_name
                  ? member.member.first_name[0]
                  : member.member?.email?.[0]?.toUpperCase() || "?"}
              </div>
              <div>
                <div class="flex items-center gap-2">
                  <p class="font-medium text-gray-900">
                    {#if member.member?.email}
                      {member.member.first_name && member.member.last_name
                        ? `${member.member.first_name} ${member.member.last_name}`
                        : member.member.name || member.member.email}
                    {:else if member.role === "OWNER"}
                      {$auth.user?.name || $auth.user?.email || "Owner"}
                    {:else}
                      Unknown
                    {/if}
                  </p>
                  <span
                    class={`inline-flex items-center gap-1 px-2 py-1 text-xs font-medium rounded-full ${getRoleBadgeClass(member.role || "")}`}
                  >
                    <RoleIcon class="h-3 w-3" />
                    {member.role?.replace("Role", "") || "Unknown"}
                  </span>
                </div>
                <p class="text-sm text-gray-500">
                  {#if member.member?.email}
                    {member.member.email}
                  {:else if member.role === "OWNER"}
                    {$auth.user?.email || ""}
                  {/if}
                </p>
              </div>
            </div>

            <div class="flex items-center gap-4">
              <div class="text-right">
                <p class="text-xs text-gray-500">
                  {member.accepted_at ? "Joined" : "Invited"}
                </p>
                <p class="text-xs font-medium text-gray-700">
                  {formatDate(member.accepted_at || member.invited_at || "")}
                </p>
              </div>

              <RoleGate roles={[ROLES.owner, ROLES.admin]}>
                <div class="flex items-center gap-2">
                  {#if member.role !== ROLES.owner}
                    <select
                      class="text-sm border-gray-300 rounded-md"
                      value={member.role}
                      onchange={(e) =>
                        updateRole(
                          member.id || "",
                          e?.currentTarget?.value as Role,
                        )}
                    >
                      <option value={ROLES.admin}>Admin</option>
                      <option value={ROLES.manager}>Manager</option>
                      <option value={ROLES.viewer}>Viewer</option>
                    </select>
                  {/if}

                  {#if member.role !== ROLES.owner}
                    <Button
                      variant="ghost"
                      size="sm"
                      onclick={() => removeMember(member.id || "")}
                      disabled={removingMemberId === member.id}
                    >
                      {#if removingMemberId === member.id}
                        <Loader2 class="h-4 w-4 animate-spin" />
                      {:else}
                        <Trash2 class="h-4 w-4 text-red-500" />
                      {/if}
                    </Button>
                  {/if}
                </div>
              </RoleGate>
            </div>
          </div>

          {#if !member.accepted_at}
            <div
              class="mt-3 flex items-center justify-between text-sm text-amber-600 bg-amber-50 px-3 py-2 rounded-md"
            >
              <div class="flex items-center gap-2">
                <Mail class="h-4 w-4" />
                Invitation pending
              </div>
              <RoleGate roles={[ROLES.owner, ROLES.admin]}>
                <Button
                  variant="ghost"
                  size="sm"
                  onclick={() => resendInvitation(member.id || "")}
                  disabled={resendingInvitationId === member.id}
                >
                  {#if resendingInvitationId === member.id}
                    <Loader2 class="h-3 w-3 animate-spin" />
                  {:else}
                    <Mail class="h-3 w-3 mr-1" />
                    Resend
                  {/if}
                </Button>
              </RoleGate>
            </div>
          {/if}
        </div>
      {/each}
    </div>

    <!-- Role Descriptions -->
    <div class="mt-8 bg-gray-50 rounded-lg p-6">
      <h3 class="text-sm font-semibold text-gray-900 mb-4">Role Permissions</h3>
      <div class="space-y-3 text-sm">
        <div class="flex items-start gap-3">
          <Shield class="h-5 w-5 text-purple-500 mt-0.5" />
          <div>
            <p class="font-medium text-gray-900">Owner</p>
            <p class="text-gray-600">Full access to all features.</p>
          </div>
        </div>
        <div class="flex items-start gap-3">
          <Settings class="h-5 w-5 text-blue-500 mt-0.5" />
          <div>
            <p class="font-medium text-gray-900">Admin</p>
            <p class="text-gray-600">
              Full access with exceptions in team management.
            </p>
          </div>
        </div>
        <div class="flex items-start gap-3">
          <Users class="h-5 w-5 text-green-500 mt-0.5" />
          <div>
            <p class="font-medium text-gray-900">Manager</p>
            <p class="text-gray-600">
              Read and write rights to feedback content.
            </p>
          </div>
        </div>
        <div class="flex items-start gap-3">
          <Eye class="h-5 w-5 text-gray-500 mt-0.5" />
          <div>
            <p class="font-medium text-gray-900">Viewer</p>
            <p class="text-gray-600">Read rights to feedback content.</p>
          </div>
        </div>
      </div>
    </div>
  {/if}

  <!-- Invite Modal -->
  <Modal bind:open={showInviteModal} title="Invite Team Member">
    <form
      onsubmit={(e) => {
        e.preventDefault();
        inviteMember();
      }}
      class="space-y-4"
    >
      <div>
        <label
          for="invite-email"
          class="block text-sm font-medium text-gray-700 mb-1"
        >
          Email Address
        </label>
        <Input
          id="invite-email"
          type="email"
          bind:value={inviteEmail}
          placeholder="colleague@example.com"
          required
        />
      </div>

      <div>
        <label
          for="invite-role"
          class="block text-sm font-medium text-gray-700 mb-1"
        >
          Role
        </label>
        <select
          id="invite-role"
          bind:value={inviteRole}
          class="w-full border-gray-300 rounded-md"
        >
          <option value="ADMIN">Admin</option>
          <option value="MANAGER">Manager</option>
          <option value="VIEWER">Viewer</option>
        </select>
      </div>

      <div class="flex justify-end gap-2 pt-4">
        <Button
          type="button"
          variant="outline"
          onclick={() => (showInviteModal = false)}
          disabled={inviting}
        >
          Cancel
        </Button>
        <Button
          type="submit"
          variant="gradient"
          disabled={inviting || !inviteEmail}
        >
          {#if inviting}
            <Loader2 class="h-4 w-4 mr-2 animate-spin" />
            Sending...
          {:else}
            <Mail class="h-4 w-4 mr-2" />
            Send Invitation
          {/if}
        </Button>
      </div>
    </form>
  </Modal>
</div>
