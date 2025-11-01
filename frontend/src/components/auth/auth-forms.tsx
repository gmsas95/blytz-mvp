'use client';

import { Loader2, Eye, EyeOff } from 'lucide-react';
import { useRouter } from 'next/navigation';
import { useState } from 'react';

import { Alert, AlertDescription } from '@/components/ui/alert';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { useAuth } from '@/contexts/auth-context';
import { authRateLimiter } from '@/lib/auth-middleware';

interface LoginFormProps {
  onSuccess?: () => void;
}

export function LoginForm({ onSuccess }: LoginFormProps) {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [showPassword, setShowPassword] = useState(false);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const { login } = useAuth();
  const router = useRouter();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsLoading(true);
    setError(null);

    // Rate limiting check
    const isRateLimited = authRateLimiter.isRateLimited(email);
    if (isRateLimited) {
      const resetTime = authRateLimiter.getResetTime(email);
      setError(`Too many login attempts. Try again after ${resetTime?.toLocaleTimeString()}`);
      setIsLoading(false);
      return;
    }

    try {
      const result = await login(email, password);
      if (result.success) {
        onSuccess?.();
        router.push('/');
        router.refresh();
      } else {
        setError(result.error || 'Login failed');
      }
    } catch (err) {
      setError('An unexpected error occurred');
    } finally {
      setIsLoading(false);
    }
  };

  const remainingAttempts = authRateLimiter.getRemainingAttempts(email);

  return (
    <form onSubmit={handleSubmit} className="space-y-4" noValidate>
      {error && (
        <Alert variant="destructive" role="alert" aria-live="assertive">
          <AlertDescription>{error}</AlertDescription>
        </Alert>
      )}

      <div className="space-y-2">
        <Label htmlFor="email">Email address</Label>
        <Input
          id="email"
          type="email"
          placeholder="Enter your email address"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          required
          disabled={isLoading}
          autoComplete="email"
          aria-describedby="email-help"
          aria-invalid={!!error && error.includes('email')}
        />
        <div id="email-help" className="text-sm text-muted-foreground">
          Enter the email address you used to register
        </div>
      </div>

      <div className="space-y-2">
        <Label htmlFor="password">Password</Label>
        <div className="relative">
          <Input
            id="password"
            type={showPassword ? 'text' : 'password'}
            placeholder="Enter your password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
            disabled={isLoading}
            autoComplete="current-password"
            aria-describedby="password-help toggle-password-visibility"
            aria-invalid={!!error && error.includes('password')}
          />
          <Button
            id="toggle-password-visibility"
            type="button"
            variant="ghost"
            size="sm"
            className="absolute right-0 top-0 h-full px-3 py-2 hover:bg-transparent"
            onClick={() => setShowPassword(!showPassword)}
            disabled={isLoading}
            aria-label={showPassword ? 'Hide password' : 'Show password'}
            aria-pressed={showPassword}
          >
            {showPassword ? <EyeOff className="h-4 w-4" /> : <Eye className="h-4 w-4" />}
          </Button>
        </div>
        <div id="password-help" className="text-sm text-muted-foreground">
          Enter your account password
        </div>
      </div>

      {remainingAttempts < 5 && (
        <p className="text-sm text-muted-foreground" role="status" aria-live="polite">
          {remainingAttempts} login attempts remaining
        </p>
      )}

      <Button type="submit" className="w-full" disabled={isLoading} aria-describedby="login-status">
        {isLoading && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
        Sign in
      </Button>

      {isLoading && (
        <div id="login-status" className="sr-only" role="status" aria-live="polite">
          Signing in to your account
        </div>
      )}

      <div
        className="text-sm text-muted-foreground bg-muted/30 rounded-lg p-3"
        role="status"
        aria-live="polite"
      >
        <p className="font-medium mb-1">Demo credentials:</p>
        <p>Email: demo@blytz.app</p>
        <p>Password: demo123</p>
      </div>
    </form>
  );
}

interface RegisterFormProps {
  onSuccess?: () => void;
  onSwitchToLogin?: () => void;
}

export function RegisterForm({ onSuccess, onSwitchToLogin }: RegisterFormProps) {
  const [formData, setFormData] = useState({
    name: '',
    displayName: '',
    email: '',
    password: '',
    confirmPassword: '',
  });
  const [showPassword, setShowPassword] = useState(false);
  const [showConfirmPassword, setShowConfirmPassword] = useState(false);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const { register } = useAuth();
  const router = useRouter();

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setFormData((prev) => ({
      ...prev,
      [e.target.name]: e.target.value,
    }));
  };

  const validateForm = () => {
    if (!formData.name.trim()) {
      setError('Name is required');
      return false;
    }

    if (!formData.email.trim() || !formData.email.includes('@')) {
      setError('Valid email is required');
      return false;
    }

    if (formData.password.length < 6) {
      setError('Password must be at least 6 characters');
      return false;
    }

    if (formData.password !== formData.confirmPassword) {
      setError('Passwords do not match');
      return false;
    }

    return true;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);

    if (!validateForm()) return;

    setIsLoading(true);

    try {
      const result = await register({
        name: formData.name,
        displayName: formData.displayName || formData.name,
        email: formData.email,
        password: formData.password,
      });

      if (result.success) {
        onSuccess?.();
        router.push('/');
        router.refresh();
      } else {
        setError(result.error || 'Registration failed');
      }
    } catch (err) {
      setError('An unexpected error occurred');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-4" noValidate>
      {error && (
        <Alert variant="destructive" role="alert" aria-live="assertive">
          <AlertDescription>{error}</AlertDescription>
        </Alert>
      )}

      <div className="space-y-2">
        <Label htmlFor="name">Full Name</Label>
        <Input
          id="name"
          name="name"
          type="text"
          placeholder="Enter your full name"
          value={formData.name}
          onChange={handleChange}
          required
          disabled={isLoading}
          autoComplete="name"
          aria-describedby="name-help"
          aria-invalid={!!error && error.includes('name')}
        />
        <div id="name-help" className="text-sm text-muted-foreground">
          Enter your legal full name
        </div>
      </div>

      <div className="space-y-2">
        <Label htmlFor="displayName">Display Name (optional)</Label>
        <Input
          id="displayName"
          name="displayName"
          type="text"
          placeholder="How others will see you"
          value={formData.displayName}
          onChange={handleChange}
          disabled={isLoading}
          autoComplete="nickname"
          aria-describedby="displayname-help"
        />
        <div id="displayname-help" className="text-sm text-muted-foreground">
          This is how other users will see you on the platform
        </div>
      </div>

      <div className="space-y-2">
        <Label htmlFor="register-email">Email address</Label>
        <Input
          id="register-email"
          name="email"
          type="email"
          placeholder="Enter your email address"
          value={formData.email}
          onChange={handleChange}
          required
          disabled={isLoading}
          autoComplete="email"
          aria-describedby="register-email-help"
          aria-invalid={!!error && error.includes('email')}
        />
        <div id="register-email-help" className="text-sm text-muted-foreground">
          We'll use this for account notifications
        </div>
      </div>

      <div className="space-y-2">
        <Label htmlFor="register-password">Password</Label>
        <div className="relative">
          <Input
            id="register-password"
            name="password"
            type={showPassword ? 'text' : 'password'}
            placeholder="Create a password (min. 6 characters)"
            value={formData.password}
            onChange={handleChange}
            required
            disabled={isLoading}
            autoComplete="new-password"
            aria-describedby="register-password-help toggle-register-password-visibility"
            aria-invalid={!!error && error.includes('password')}
          />
          <Button
            id="toggle-register-password-visibility"
            type="button"
            variant="ghost"
            size="sm"
            className="absolute right-0 top-0 h-full px-3 py-2 hover:bg-transparent"
            onClick={() => setShowPassword(!showPassword)}
            disabled={isLoading}
            aria-label={showPassword ? 'Hide password' : 'Show password'}
            aria-pressed={showPassword}
          >
            {showPassword ? <EyeOff className="h-4 w-4" /> : <Eye className="h-4 w-4" />}
          </Button>
        </div>
        <div id="register-password-help" className="text-sm text-muted-foreground">
          Must be at least 6 characters long
        </div>
      </div>

      <div className="space-y-2">
        <Label htmlFor="confirmPassword">Confirm Password</Label>
        <div className="relative">
          <Input
            id="confirmPassword"
            name="confirmPassword"
            type={showConfirmPassword ? 'text' : 'password'}
            placeholder="Confirm your password"
            value={formData.confirmPassword}
            onChange={handleChange}
            required
            disabled={isLoading}
            autoComplete="new-password"
            aria-describedby="confirm-password-help toggle-confirm-password-visibility"
            aria-invalid={!!error && error.includes('match')}
          />
          <Button
            id="toggle-confirm-password-visibility"
            type="button"
            variant="ghost"
            size="sm"
            className="absolute right-0 top-0 h-full px-3 py-2 hover:bg-transparent"
            onClick={() => setShowConfirmPassword(!showConfirmPassword)}
            disabled={isLoading}
            aria-label={showConfirmPassword ? 'Hide confirm password' : 'Show confirm password'}
            aria-pressed={showConfirmPassword}
          >
            {showConfirmPassword ? <EyeOff className="h-4 w-4" /> : <Eye className="h-4 w-4" />}
          </Button>
        </div>
        <div id="confirm-password-help" className="text-sm text-muted-foreground">
          Re-enter your password to confirm
        </div>
      </div>

      <Button
        type="submit"
        className="w-full"
        disabled={isLoading}
        aria-describedby="register-status"
      >
        {isLoading && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
        Create account
      </Button>

      {isLoading && (
        <div id="register-status" className="sr-only" role="status" aria-live="polite">
          Creating your account
        </div>
      )}

      <div className="text-center">
        <p className="text-sm text-muted-foreground">
          Already have an account?{' '}
          <button
            type="button"
            onClick={onSwitchToLogin}
            className="text-primary hover:underline focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2 rounded px-1"
            disabled={isLoading}
          >
            Sign in
          </button>
        </p>
      </div>
    </form>
  );
}
