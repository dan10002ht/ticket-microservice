import { getUserRepository, getRoleRepository } from '../../repositories/repositoryFactory.js';

// Get repository instances from factory
const userRepository = getUserRepository();
const roleRepository = getRoleRepository();

/**
 * Get user roles
 */
export async function getUserRoles(userId) {
  try {
    const user = await userRepository.findByPublicId(userId);
    if (!user) {
      throw new Error('User not found');
    }

    const roles = await roleRepository.findByUserId(user.id);
    return roles.map((role) => ({
      id: role.public_id,
      name: role.name,
      description: role.description,
    }));
  } catch (error) {
    throw new Error(`Failed to get user roles: ${error.message}`);
  }
}

/**
 * Assign role to user
 */
export async function assignRoleToUser(userId, roleId) {
  try {
    // Check if user exists
    const user = await userRepository.findByPublicId(userId);
    if (!user) {
      throw new Error('User not found');
    }

    // Check if role exists
    const role = await roleRepository.findByPublicId(roleId);
    if (!role) {
      throw new Error('Role not found');
    }

    // Check if user already has this role
    const existingRoles = await roleRepository.findByUserId(user.id);
    const hasRole = existingRoles.some((r) => r.id === role.id);
    if (hasRole) {
      throw new Error('User already has this role');
    }

    // Assign role using internal ids
    await roleRepository.assignToUser(user.id, role.id);

    return {
      message: 'Role assigned successfully',
      role: {
        id: role.public_id,
        name: role.name,
        description: role.description,
      },
    };
  } catch (error) {
    throw new Error(`Failed to assign role: ${error.message}`);
  }
}

/**
 * Remove role from user
 */
export async function removeRoleFromUser(userId, roleId) {
  try {
    // Get user and role by public ids
    const user = await userRepository.findByPublicId(userId);
    if (!user) {
      throw new Error('User not found');
    }

    const role = await roleRepository.findByPublicId(roleId);
    if (!role) {
      throw new Error('Role not found');
    }

    // Check if user has this role
    const existingRoles = await roleRepository.findByUserId(user.id);
    const hasRole = existingRoles.some((r) => r.id === role.id);
    if (!hasRole) {
      throw new Error('User does not have this role');
    }

    // Remove role using internal ids
    await roleRepository.removeFromUser(user.id, role.id);

    return {
      message: 'Role removed successfully',
    };
  } catch (error) {
    throw new Error(`Failed to remove role: ${error.message}`);
  }
}

/**
 * Get all roles
 */
export async function getAllRoles() {
  try {
    const roles = await roleRepository.findAll();
    return roles.map((role) => ({
      id: role.public_id,
      name: role.name,
      description: role.description,
    }));
  } catch (error) {
    throw new Error(`Failed to get all roles: ${error.message}`);
  }
}

/**
 * Get role by ID
 */
export async function getRoleById(roleId) {
  try {
    const role = await roleRepository.findByPublicId(roleId);
    if (!role) {
      throw new Error('Role not found');
    }

    return {
      id: role.public_id,
      name: role.name,
      description: role.description,
    };
  } catch (error) {
    throw new Error(`Failed to get role: ${error.message}`);
  }
}

/**
 * Check if user has specific role
 */
export async function userHasRole(userId, roleName) {
  try {
    const user = await userRepository.findByPublicId(userId);
    if (!user) {
      return false;
    }

    const roles = await roleRepository.findByUserId(user.id);
    return roles.some((role) => role.name === roleName);
  } catch (error) {
    throw new Error(`Failed to check user role: ${error.message}`);
  }
}
