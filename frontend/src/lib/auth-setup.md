# Authentication System Documentation

## Overview

This document describes the complete JWT-based authentication system implemented
for the Blytz frontend. The system integrates with the backend authentication
service at `http://localhost:8084` and provides secure user management
throughout the application.

## Architecture

### Core Components

1. **AuthContext** (`/contexts/auth-context.tsx`)
   - Centralized authentication state management
   - JWT token handling and storage
   - User session management
   - Automatic token validation

2. **Protected Routes** (`/components/auth/protected-route.tsx`)
   - HOC for protecting authenticated routes
   - Automatic redirect logic
   - Loading states during auth checks

3. **Auth Forms** (`/components/auth/auth-forms.tsx`)
   - Login and register form components
   - Rate limiting protection
   - Comprehensive error handling
   - Form validation

4. **Enhanced API Adapter** (`/lib/api-adapter.ts`)
   - Automatic JWT token injection
   - Token refresh logic
   - Error handling for unauthorized requests

5. **Auth Middleware** (`/lib/auth-middleware.ts`)
   - Rate limiting utilities
   - Permission checking helpers
   - Token management utilities

## Authentication Flow

### Login Process

1. User enters credentials in login form
2. Form validation and rate limiting check
3. API call to `/auth/login` endpoint
4. Backend validates credentials and returns JWT token
5. Token stored securely in localStorage (httpOnly cookies in production)
6. User state updated in AuthContext
7. Redirect to dashboard/home page

### Registration Process

1. User fills registration form
2. Client-side validation
3. API call to `/auth/register` endpoint
4. Backend creates user account
5. Auto-login with new credentials
6. Token stored and user state updated

### Token Management

- **Storage**: localStorage for development, httpOnly cookies for production
- **Validation**: Token validated on app initialization and API calls
- **Expiration**: Automatic logout on token expiration
- **Refresh**: Token refresh handled automatically

## Usage Examples

### Protecting Routes

```tsx
import { ProtectedRoute } from '@/components/auth/protected-route';

export default function ProtectedPage() {
  return (
    <ProtectedRoute>
      <div>This page requires authentication</div>
    </ProtectedRoute>
  );
}
```

### Using Auth Context

```tsx
import { useAuth } from '@/contexts/auth-context';

export default function MyComponent() {
  const { user, isAuthenticated, login, logout } = useAuth();

  const handleLogin = async () => {
    const result = await login('user@example.com', 'password');
    if (result.success) {
      // Login successful
    } else {
      // Handle error
      console.error(result.error);
    }
  };

  return (
    <div>
      {isAuthenticated ? (
        <div>Welcome, {user?.name}!</div>
      ) : (
        <button onClick={handleLogin}>Login</button>
      )}
    </div>
  );
}
```

### API Calls with Authentication

The API adapter automatically includes JWT tokens in requests:

```tsx
import { api } from '@/lib/api-adapter';

// This call will automatically include the auth token
const response = await api.getCart();

// Protected endpoints work seamlessly
const userResponse = await api.getCurrentUser();
```

## Security Features

### Rate Limiting

- Login attempts limited to 5 attempts per 15 minutes
- Tracking based on email address
- Automatic lockout with reset time display

### Token Security

- Secure token storage (httpOnly cookies in production)
- Automatic token validation
- Proper token cleanup on logout
- XSS protection considerations

### Input Validation

- Client-side form validation
- Server-side validation required
- SQL injection protection (backend)
- XSS prevention

## Configuration

### Environment Variables

```env
# Backend API URL
NEXT_PUBLIC_API_URL=http://localhost:8084

# Mode (mock for development, remote for production)
MODE=mock

# Fiuu payment sandbox mode
NEXT_PUBLIC_FIUU_SANDBOX=true
```

### Backend Integration

The frontend expects the following API endpoints:

- `POST /auth/login` - User authentication
- `POST /auth/register` - User registration
- `POST /auth/logout` - User logout
- `GET /auth/me` - Get current user

Expected response format:

```json
{
  "success": true,
  "data": {
    "user": {
      "id": "user_id",
      "name": "User Name",
      "email": "user@example.com",
      "isSeller": false,
      "rating": 4.5,
      "totalSales": 1000
    },
    "token": "jwt_token_here"
  }
}
```

## Development vs Production

### Development Mode (Mock)

- Uses mock data for development
- Demo credentials: demo@blytz.app / demo123
- No actual network requests
- Simulated authentication flow

### Production Mode (Remote)

- Real API calls to backend
- JWT tokens from authentication service
- Secure httpOnly cookies
- Full authentication integration

## Error Handling

### Authentication Errors

- Invalid credentials
- Token expiration
- Network errors
- Rate limiting

### User Experience

- Loading states during auth operations
- Clear error messages
- Automatic redirect on auth failure
- Graceful degradation

## Testing

### Mock Authentication

```typescript
// In development mode
const result = await api.login('demo@blytz.app', 'demo123');
// Returns: { success: true, data: mockUser }
```

### Protected Route Testing

```typescript
// Protected routes automatically redirect unauthenticated users
// Accessing /profile without authentication redirects to /auth
```

## Future Enhancements

1. **Social Login**: Google, Facebook, Apple integration
2. **Two-Factor Authentication**: SMS/Email verification
3. **Password Reset**: Email-based password recovery
4. **Session Management**: Multiple device session handling
5. **Admin Panel**: User management interface
6. **Audit Logging**: Track authentication events

## Troubleshooting

### Common Issues

1. **Token Not Persisting**: Check localStorage/clearance
2. **CORS Issues**: Verify backend CORS configuration
3. **Invalid Token**: Check token format and expiration
4. **Route Protection**: Ensure ProtectedRoute wrapper is used
5. **API Errors**: Check network connectivity and backend status

### Debug Tools

- Browser DevTools localStorage inspection
- Network tab for API request debugging
- React DevTools for state inspection
- Console logs for authentication events

## Files Created/Modified

### New Files

- `/src/contexts/auth-context.tsx` - Authentication context
- `/src/components/auth/protected-route.tsx` - Route protection
- `/src/components/auth/auth-forms.tsx` - Login/Register forms
- `/src/lib/auth-middleware.ts` - Authentication utilities
- `/src/components/ui/tabs.tsx` - UI component
- `/src/components/ui/dropdown-menu.tsx` - UI component
- `/src/components/ui/avatar.tsx` - UI component

### Modified Files

- `/src/app/layout.tsx` - Added AuthProvider
- `/src/app/auth/page.tsx` - Complete auth page
- `/src/app/profile/page.tsx` - Protected profile page
- `/src/components/layout/header.tsx` - User state integration
- `/src/lib/api-adapter.ts` - JWT token injection

## Summary

The authentication system provides:

✅ **Secure JWT-based authentication** ✅ **Protected routes with automatic
redirects** ✅ **User state management throughout the app** ✅ **Rate limiting
protection** ✅ **Comprehensive error handling** ✅ **Development and production
modes** ✅ **Accessible UI components** ✅ **Mobile-responsive design** ✅
**TypeScript support** ✅ **Backend integration ready**

The system is production-ready and integrates seamlessly with the existing Blytz
application architecture.
