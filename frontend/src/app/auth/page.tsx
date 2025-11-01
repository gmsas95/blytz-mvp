'use client';

import { Loader2 } from 'lucide-react';
import dynamic from 'next/dynamic';
import { useState, Suspense } from 'react';

import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';


// Dynamic imports for auth forms
const LoginForm = dynamic(() => import('@/components/auth/auth-forms').then(mod => ({ default: mod.LoginForm })), {
  loading: () => <div className="flex items-center justify-center py-8"><Loader2 className="w-6 h-6 animate-spin" /></div>,
  ssr: false
});

const RegisterForm = dynamic(() => import('@/components/auth/auth-forms').then(mod => ({ default: mod.RegisterForm })), {
  loading: () => <div className="flex items-center justify-center py-8"><Loader2 className="w-6 h-6 animate-spin" /></div>,
  ssr: false
});

const ProtectedRoute = dynamic(() => import('@/components/auth/protected-route').then(mod => ({ default: mod.ProtectedRoute })), {
  ssr: false
});

export default function AuthPage() {
  const [activeTab, setActiveTab] = useState('login');

  return (
    <Suspense fallback={
      <section className="w-full py-16 md:py-24">
        <div className="container mx-auto px-4 max-w-md">
          <div className="flex items-center justify-center py-16">
            <Loader2 className="w-8 h-8 animate-spin" />
          </div>
        </div>
      </section>
    }>
      <ProtectedRoute requireAuth={false}>
        <section className="w-full py-16 md:py-24">
          <div className="container mx-auto px-4 max-w-md">
            <Card>
              <CardHeader className="text-center">
                <CardTitle className="text-2xl md:text-3xl font-bold">Welcome to Blytz</CardTitle>
                <CardDescription>
                  Sign in to your account or create a new one to start bidding on amazing products
                </CardDescription>
              </CardHeader>
              <CardContent>
                <Tabs value={activeTab} onValueChange={setActiveTab} className="w-full">
                  <TabsList className="grid w-full grid-cols-2">
                    <TabsTrigger value="login">Sign In</TabsTrigger>
                    <TabsTrigger value="register">Sign Up</TabsTrigger>
                  </TabsList>

                  <TabsContent value="login" className="mt-6">
                    <Suspense fallback={<div className="flex items-center justify-center py-8"><Loader2 className="w-6 h-6 animate-spin" /></div>}>
                      <LoginForm />
                    </Suspense>
                  </TabsContent>

                  <TabsContent value="register" className="mt-6">
                    <Suspense fallback={<div className="flex items-center justify-center py-8"><Loader2 className="w-6 h-6 animate-spin" /></div>}>
                      <RegisterForm onSwitchToLogin={() => setActiveTab('login')} />
                    </Suspense>
                  </TabsContent>
                </Tabs>
              </CardContent>
            </Card>
          </div>
        </section>
      </ProtectedRoute>
    </Suspense>
  );
}
