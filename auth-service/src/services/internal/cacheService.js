import * as redisService from '../redis/redisService.js';

// ========== USER PROFILE CACHE ==========

export async function setUserProfileCache(userId, profileData) {
  try {
    await redisService.setUserProfile(userId, profileData);
    return true;
  } catch (error) {
    return false;
  }
}

export async function getUserProfileCache(userId) {
  try {
    return await redisService.getUserProfile(userId);
  } catch (error) {
    return null;
  }
}

export async function invalidateUserProfileCache(userId) {
  try {
    await redisService.invalidateUserProfile(userId);
    return true;
  } catch (error) {
    return false;
  }
}

// ========== USER ROLES CACHE ==========

export async function setUserRolesCache(userId, rolesData) {
  try {
    await redisService.setUserRoles(userId, rolesData);
    return true;
  } catch (error) {
    return false;
  }
}

export async function getUserRolesCache(userId) {
  try {
    return await redisService.getUserRoles(userId);
  } catch (error) {
    return null;
  }
}

export async function invalidateUserRolesCache(userId) {
  try {
    await redisService.invalidateUserRoles(userId);
    return true;
  } catch (error) {
    return false;
  }
}

// ========== TOKEN VALIDATION CACHE ==========

export async function setTokenValidationCache(token, validationData) {
  try {
    await redisService.setTokenValidation(token, validationData);
    return true;
  } catch (error) {
    return false;
  }
}

export async function getTokenValidationCache(token) {
  try {
    return await redisService.getTokenValidation(token);
  } catch (error) {
    return null;
  }
}

export async function invalidateTokenValidationCache(token) {
  try {
    await redisService.invalidateTokenValidation(token);
    return true;
  } catch (error) {
    return false;
  }
}

// ========== ORGANIZATION CACHE ==========

export async function setOrganizationCache(orgId, orgData) {
  try {
    await redisService.setOrganization(orgId, orgData);
    return true;
  } catch (error) {
    return false;
  }
}

export async function getOrganizationCache(orgId) {
  try {
    return await redisService.getOrganization(orgId);
  } catch (error) {
    return null;
  }
}

export async function invalidateOrganizationCache(orgId) {
  try {
    await redisService.invalidateOrganization(orgId);
    return true;
  } catch (error) {
    return false;
  }
}

// ========== ROLE PERMISSIONS CACHE ==========

export async function setRolePermissionsCache(roleId, permissionsData) {
  try {
    await redisService.setRolePermissions(roleId, permissionsData);
    return true;
  } catch (error) {
    return false;
  }
}

export async function getRolePermissionsCache(roleId) {
  try {
    return await redisService.getRolePermissions(roleId);
  } catch (error) {
    return null;
  }
}

export async function invalidateRolePermissionsCache(roleId) {
  try {
    await redisService.invalidateRolePermissions(roleId);
    return true;
  } catch (error) {
    return false;
  }
}

// ========== SESSION CACHE ==========

export async function setSessionCache(sessionId, sessionData) {
  try {
    await redisService.setSession(sessionId, sessionData);
    return true;
  } catch (error) {
    return false;
  }
}

export async function getSessionCache(sessionId) {
  try {
    return await redisService.getSession(sessionId);
  } catch (error) {
    return null;
  }
}

export async function deleteSessionCache(sessionId) {
  try {
    await redisService.deleteSession(sessionId);
    return true;
  } catch (error) {
    return false;
  }
}

export async function setUserSessionsCache(userId, sessionsData) {
  try {
    await redisService.setUserSessions(userId, sessionsData);
    return true;
  } catch (error) {
    return false;
  }
}

export async function getUserSessionsCache(userId) {
  try {
    return await redisService.getUserSessions(userId);
  } catch (error) {
    return null;
  }
}

export async function invalidateUserSessionsCache(userId) {
  try {
    await redisService.invalidateUserSessions(userId);
    return true;
  } catch (error) {
    return false;
  }
}

// ========== HEALTH CHECK ==========

export async function cacheHealthCheck() {
  return await redisService.healthCheck();
}

// ========== UTILITY METHODS ==========

export async function invalidateUserCache(userId) {
  try {
    await Promise.all([
      redisService.invalidateUserProfile(userId),
      redisService.invalidateUserRoles(userId),
      redisService.invalidateUserSessions(userId),
    ]);
    return true;
  } catch (error) {
    return false;
  }
}

export async function invalidateTokenCache(token) {
  try {
    await redisService.invalidateTokenValidation(token);
    return true;
  } catch (error) {
    return false;
  }
}

// ========== LEGACY METHODS (for backward compatibility) ==========

export async function cacheUserProfile(userId, userData) {
  return await setUserProfileCache(userId, userData);
}

export async function getCachedUserProfile(userId) {
  return await getUserProfileCache(userId);
}

export async function cacheUserRoles(userId, roles) {
  return await setUserRolesCache(userId, roles);
}

export async function getCachedUserRoles(userId) {
  return await getUserRolesCache(userId);
}

export async function cacheTokenValidation(token, validationData) {
  return await setTokenValidationCache(token, validationData);
}

export async function getCachedTokenValidation(token) {
  return await getTokenValidationCache(token);
}

export async function cacheOrganization(orgId, orgData) {
  return await setOrganizationCache(orgId, orgData);
}

export async function getCachedOrganization(orgId) {
  return await getOrganizationCache(orgId);
}

export async function cacheSession(sessionId, sessionData) {
  return await setSessionCache(sessionId, sessionData);
}

export async function getCachedSession(sessionId) {
  return await getSessionCache(sessionId);
}
