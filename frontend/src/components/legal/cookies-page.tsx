'use client';

import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';

export function CookiesPage() {
  return (
    <div className="container mx-auto px-4 py-8 max-w-4xl">
      <div className="mb-8">
        <h1 className="text-4xl font-bold mb-4">Cookie Policy</h1>
        <p className="text-lg text-muted-foreground">
          How Blytz uses cookies and similar technologies to enhance your experience.
        </p>
      </div>

      <div className="space-y-6">
        <Card>
          <CardHeader>
            <CardTitle>What Are Cookies</CardTitle>
          </CardHeader>
          <CardContent>
            <p>
              Cookies are small text files stored on your device when you visit websites. 
              They help us remember your preferences and improve your experience.
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>How We Use Cookies</CardTitle>
          </CardHeader>
          <CardContent>
            <p>We use cookies to:</p>
            <ul className="list-disc list-inside mt-2 space-y-1">
              <li>Keep you logged in to your account</li>
              <li>Remember your preferences and settings</li>
              <li>Understand how you use our service</li>
              <li>Personalize content and advertisements</li>
              <li>Analyze website performance</li>
            </ul>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Types of Cookies</CardTitle>
          </CardHeader>
          <CardContent>
            <p><strong>Essential Cookies:</strong> Required for basic site functionality</p>
            <p><strong>Performance Cookies:</strong> Help us understand site usage</p>
            <p><strong>Functional Cookies:</strong> Remember your preferences</p>
            <p><strong>Marketing Cookies:</strong> Used for advertising purposes</p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Managing Cookies</CardTitle>
          </CardHeader>
          <CardContent>
            <p>
              You can control cookies through your browser settings. However, 
              disabling cookies may affect your ability to use some features of our service.
            </p>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}