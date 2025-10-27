export const environment = {
    production: false,
    // All requests now go through API Gateway
    apiGateway: 'http://localhost:8080',
    // Individual service hosts for reference (no longer used directly)
    apiHost: 'http://localhost:8080',
    authApiHost: 'http://localhost:8080/api/v1/auth',
    blogApiHost: 'http://localhost:8080/api/v1/blogs',
    stakeholdersApiHost: 'http://localhost:8080/api/v1',
    tourApiHost: 'http://localhost:8080/api/v1/tours',
    followerApiHost: 'http://localhost:8080/api/followers',
    purchaseApiHost: 'http://localhost:8080/api/v1/cart'
  };
  