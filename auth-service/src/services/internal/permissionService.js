import {
  getUserRepository,
  getRoleRepository,
  getPermissionRepository,
} from '../../repositories/repositoryFactory.js';
import { sanitizeUserForResponse } from '../../utils/sanitizers.js';
import { getCachedPermission, setCachedPermission } from './cacheService.js';

// Get repository instances from factory
const userRepository = getUserRepository();
const roleRepository = getRoleRepository();
const permissionRepository = getPermissionRepository();

// ========== USER PERMISSIONS ==========

/**
 * Get user permissions
 */
export async function getUserPermissions(userId) {
  try {
    const user = await userRepository.findByPublicId(userId);
    if (!user) {
      throw new Error('User not found');
    }

    // Get user permissions from database using internal id
    const permissions = await permissionRepository.getUserPermissions(user.id);

    // Format permissions for response
    const formattedPermissions = permissions.map((permission) => ({
      id: permission.id,
      name: permission.name,
      description: permission.description,
      resource: permission.resource,
      action: permission.action,
    }));

    return {
      user: sanitizeUserForResponse(user),
      permissions: formattedPermissions,
    };
  } catch (error) {
    throw new Error(`Failed to get user permissions: ${error.message}`);
  }
}

/**
 * Check if user has specific permission
 */
export async function hasPermission(userId, permissionName) {
  try {
    // Check cache first
    const cacheKey = `perm:${userId}:${permissionName}`;
    const cachedResult = await getCachedPermission(cacheKey);
    if (cachedResult !== null) {
      return cachedResult;
    }

    const user = await userRepository.findByPublicId(userId);
    if (!user) {
      return false;
    }

    const permissions = await permissionRepository.getUserPermissions(user.id);
    const hasPermission = permissions.some((permission) => permission.name === permissionName);

    // Cache result (5 minutes)
    await setCachedPermission(cacheKey, hasPermission, 300);

    return hasPermission;
  } catch (error) {
    throw new Error(`Failed to check permission: ${error.message}`);
  }
}

/**
 * Check if user has permission for resource and action
 */
export async function hasResourcePermission(userId, resource, action) {
  try {
    // Check cache first
    const cacheKey = `resource_perm:${userId}:${resource}:${action}`;
    const cachedResult = await getCachedPermission(cacheKey);
    if (cachedResult !== null) {
      return cachedResult;
    }

    const user = await userRepository.findByPublicId(userId);
    if (!user) {
      return false;
    }

    const permissions = await permissionRepository.getUserPermissions(user.id);
    const hasPermission = permissions.some(
      (permission) => permission.resource === resource && permission.action === action
    );

    // Cache result (5 minutes)
    await setCachedPermission(cacheKey, hasPermission, 300);

    return hasPermission;
  } catch (error) {
    throw new Error(`Failed to check resource permission: ${error.message}`);
  }
}

/**
 * Batch check multiple permissions for a user
 */
export async function batchCheckPermissions(userId, permissionNames) {
  try {
    const results = {};

    // Check cache first for each permission
    const cachePromises = permissionNames.map(async (permissionName) => {
      const cacheKey = `perm:${userId}:${permissionName}`;
      const cachedResult = await getCachedPermission(cacheKey);
      if (cachedResult !== null) {
        results[permissionName] = cachedResult;
        return { permissionName, cached: true, result: cachedResult };
      }
      return { permissionName, cached: false };
    });

    const cacheResults = await Promise.all(cachePromises);
    const uncachedPermissions = cacheResults.filter((r) => !r.cached).map((r) => r.permissionName);

    // If all permissions are cached, return results
    if (uncachedPermissions.length === 0) {
      return results;
    }

    // Get user permissions from database for uncached permissions
    const user = await userRepository.findByPublicId(userId);
    if (!user) {
      // Return false for all uncached permissions
      uncachedPermissions.forEach((permissionName) => {
        results[permissionName] = false;
      });
      return results;
    }

    const userPermissions = await permissionRepository.getUserPermissions(user.id);
    const userPermissionNames = userPermissions.map((p) => p.name);

    // Check each uncached permission
    const cacheSetPromises = uncachedPermissions.map(async (permissionName) => {
      const hasPermission = userPermissionNames.includes(permissionName);
      results[permissionName] = hasPermission;

      // Cache result
      const cacheKey = `perm:${userId}:${permissionName}`;
      await setCachedPermission(cacheKey, hasPermission, 300);
    });

    await Promise.all(cacheSetPromises);

    return results;
  } catch (error) {
    throw new Error(`Failed to batch check permissions: ${error.message}`);
  }
}

// ========== ROLE MANAGEMENT ==========

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
      id: role.id,
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

// ========== PERMISSION MANAGEMENT ==========

/**
 * Get all permissions
 */
export async function getAllPermissions() {
  try {
    const permissions = await permissionRepository.findAll();
    return permissions.map((permission) => ({
      id: permission.id,
      name: permission.name,
      description: permission.description,
      resource: permission.resource,
      action: permission.action,
    }));
  } catch (error) {
    throw new Error(`Failed to get permissions: ${error.message}`);
  }
}

/**
 * Get permissions by resource
 */
export async function getPermissionsByResource(resource) {
  try {
    const permissions = await permissionRepository.findByResource(resource);
    return permissions.map((permission) => ({
      id: permission.id,
      name: permission.name,
      description: permission.description,
      resource: permission.resource,
      action: permission.action,
    }));
  } catch (error) {
    throw new Error(`Failed to get permissions by resource: ${error.message}`);
  }
}

/**
 * Create new permission
 */
export async function createPermission(permissionData) {
  try {
    const permission = await permissionRepository.create(permissionData);
    return {
      id: permission.id,
      name: permission.name,
      description: permission.description,
      resource: permission.resource,
      action: permission.action,
    };
  } catch (error) {
    throw new Error(`Failed to create permission: ${error.message}`);
  }
}

/**
 * Assign permission to role
 */
export async function assignPermissionToRole(roleId, permissionId) {
  try {
    await permissionRepository.assignToRole(roleId, permissionId);
    return {
      message: 'Permission assigned to role successfully',
    };
  } catch (error) {
    throw new Error(`Failed to assign permission to role: ${error.message}`);
  }
}

/**
 * Remove permission from role
 */
export async function removePermissionFromRole(roleId, permissionId) {
  try {
    await permissionRepository.removeFromRole(roleId, permissionId);
    return {
      message: 'Permission removed from role successfully',
    };
  } catch (error) {
    throw new Error(`Failed to remove permission from role: ${error.message}`);
  }
}
