'use client';

import { useState } from 'react';

import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { LiveRegion } from '@/components/ui/live-region';
import { LoadingSpinner } from '@/components/ui/loading-spinner';
import { Progress } from '@/components/ui/progress-bar';
import { useAnnouncement, useFormA11y, useLiveRegion } from '@/hooks/use-accessibility';
import { announceMessage } from '@/lib/accessibility';

export default function AccessibilityTestPage() {
  const [announcement, setAnnouncement] = useState('');
  const [progress, setProgress] = useState(0);
  const [isLoading, setIsLoading] = useState(false);
  const [formData, setFormData] = useState({
    name: '',
    email: '',
    message: '',
  });
  const [formErrors, setFormErrors] = useState<Record<string, string>>({});

  const { announce } = useAnnouncement();
  const { regionRef, id: liveRegionId, announce: announceLive } = useLiveRegion('polite');
  const { formRef, announceFormError, announceFormSuccess, focusFirstError } = useFormA11y({
    onSubmit: () => {
      announce('Form submitted successfully!');
    },
    onError: (errors) => {
      console.log('Form errors:', errors);
    },
  });

  const handleAnnounce = () => {
    announceMessage(announcement, 'polite');
    setAnnouncement('');
  };

  const simulateProgress = () => {
    setProgress(0);
    const interval = setInterval(() => {
      setProgress((prev) => {
        if (prev >= 100) {
          clearInterval(interval);
          announce('Progress completed!', 'assertive');
          return 100;
        }
        return prev + 10;
      });
    }, 500);
  };

  const simulateLoading = () => {
    setIsLoading(true);
    announce('Loading started', 'polite');
    setTimeout(() => {
      setIsLoading(false);
      announce('Loading completed', 'polite');
    }, 3000);
  };

  const handleFormSubmit = (e: React.FormEvent) => {
    e.preventDefault();

    // Validate form
    const errors: Record<string, string> = {};
    if (!formData.name.trim()) errors.name = 'Name is required';
    if (!formData.email.trim()) errors.email = 'Email is required';
    else if (!formData.email.includes('@')) errors.email = 'Invalid email address';
    if (!formData.message.trim()) errors.message = 'Message is required';

    if (Object.keys(errors).length > 0) {
      setFormErrors(errors);
      announceFormError(errors);
      focusFirstError();
    } else {
      setFormErrors({});
      announceFormSuccess();
      setFormData({ name: '', email: '', message: '' });
    }
  };

  return (
    <div className="container mx-auto px-4 py-8 max-w-4xl">
      <header className="mb-12">
        <h1 className="text-4xl font-bold mb-4">Accessibility Test Page</h1>
        <p className="text-lg text-muted-foreground">
          This page demonstrates WCAG 2.1 AA compliant components and interactions.
        </p>
      </header>

      <main className="space-y-16">
        {/* Screen Reader Announcements */}
        <section aria-labelledby="announcements-heading">
          <h2 id="announcements-heading" className="text-2xl font-semibold mb-6">
            Screen Reader Announcements
          </h2>
          <div className="space-y-4 p-6 border rounded-lg">
            <div className="space-y-4">
              <Label htmlFor="announcement-input">Test Announcement</Label>
              <Input
                id="announcement-input"
                value={announcement}
                onChange={(e) => setAnnouncement(e.target.value)}
                placeholder="Enter a message to announce to screen readers"
                aria-describedby="announcement-help"
              />
              <div id="announcement-help" className="text-sm text-muted-foreground">
                This message will be announced to screen readers when you click the button
              </div>
              <Button onClick={handleAnnounce} disabled={!announcement}>
                Announce Message
              </Button>
            </div>
            <LiveRegion ref={regionRef} id={liveRegionId} />
          </div>
        </section>

        {/* Progress Indicators */}
        <section aria-labelledby="progress-heading">
          <h2 id="progress-heading" className="text-2xl font-semibold mb-6">
            Progress Indicators
          </h2>
          <div className="space-y-4 p-6 border rounded-lg">
            <Progress
              value={progress}
              label="File Upload Progress"
              showValue={true}
              className="mb-4"
            />
            <Button onClick={simulateProgress}>Simulate Progress</Button>
          </div>
        </section>

        {/* Loading States */}
        <section aria-labelledby="loading-heading">
          <h2 id="loading-heading" className="text-2xl font-semibold mb-6">
            Loading States
          </h2>
          <div className="space-y-4 p-6 border rounded-lg">
            <div className="flex items-center space-x-4">
              <LoadingSpinner size="sm" label="Small loading spinner" />
              <LoadingSpinner size="md" label="Medium loading spinner" />
              <LoadingSpinner size="lg" label="Large loading spinner" />
            </div>
            <Button onClick={simulateLoading} disabled={isLoading}>
              {isLoading ? (
                <>
                  <LoadingSpinner size="sm" label="Processing" />
                  Processing...
                </>
              ) : (
                'Simulate Loading'
              )}
            </Button>
          </div>
        </section>

        {/* Accessible Forms */}
        <section aria-labelledby="forms-heading">
          <h2 id="forms-heading" className="text-2xl font-semibold mb-6">
            Accessible Forms
          </h2>
          <div className="space-y-4 p-6 border rounded-lg">
            <form ref={formRef} onSubmit={handleFormSubmit} className="space-y-6" noValidate>
              <div className="space-y-2">
                <Label htmlFor="name">Full Name *</Label>
                <Input
                  id="name"
                  name="name"
                  type="text"
                  value={formData.name}
                  onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                  placeholder="Enter your full name"
                  aria-describedby={formErrors.name ? 'name-error' : 'name-help'}
                  aria-invalid={!!formErrors.name}
                  required
                />
                {formErrors.name ? (
                  <div id="name-error" className="text-sm text-destructive" role="alert">
                    {formErrors.name}
                  </div>
                ) : (
                  <div id="name-help" className="text-sm text-muted-foreground">
                    Enter your first and last name
                  </div>
                )}
              </div>

              <div className="space-y-2">
                <Label htmlFor="email">Email Address *</Label>
                <Input
                  id="email"
                  name="email"
                  type="email"
                  value={formData.email}
                  onChange={(e) => setFormData({ ...formData, email: e.target.value })}
                  placeholder="your.email@example.com"
                  aria-describedby={formErrors.email ? 'email-error' : 'email-help'}
                  aria-invalid={!!formErrors.email}
                  required
                />
                {formErrors.email ? (
                  <div id="email-error" className="text-sm text-destructive" role="alert">
                    {formErrors.email}
                  </div>
                ) : (
                  <div id="email-help" className="text-sm text-muted-foreground">
                    We'll use this for account notifications
                  </div>
                )}
              </div>

              <div className="space-y-2">
                <Label htmlFor="message">Message *</Label>
                <textarea
                  id="message"
                  name="message"
                  value={formData.message}
                  onChange={(e) => setFormData({ ...formData, message: e.target.value })}
                  placeholder="Enter your message here"
                  rows={4}
                  className="flex w-full rounded-2xl border border-input bg-background px-4 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50 resize-none"
                  aria-describedby={formErrors.message ? 'message-error' : 'message-help'}
                  aria-invalid={!!formErrors.message}
                  required
                />
                {formErrors.message ? (
                  <div id="message-error" className="text-sm text-destructive" role="alert">
                    {formErrors.message}
                  </div>
                ) : (
                  <div id="message-help" className="text-sm text-muted-foreground">
                    Tell us what's on your mind
                  </div>
                )}
              </div>

              <Button type="submit" className="w-full">
                Submit Form
              </Button>
            </form>
          </div>
        </section>

        {/* Keyboard Navigation */}
        <section aria-labelledby="keyboard-heading">
          <h2 id="keyboard-heading" className="text-2xl font-semibold mb-6">
            Keyboard Navigation
          </h2>
          <div className="space-y-4 p-6 border rounded-lg">
            <p className="text-muted-foreground mb-4">
              Test keyboard navigation with the following controls:
            </p>
            <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
              <Button variant="outline">Button 1</Button>
              <Button variant="outline">Button 2</Button>
              <Button variant="outline">Button 3</Button>
              <Button variant="outline">Button 4</Button>
            </div>
            <div className="text-sm text-muted-foreground">
              <p>Keyboard shortcuts to try:</p>
              <ul className="list-disc list-inside mt-2 space-y-1">
                <li>Tab/Shift+Tab: Navigate between interactive elements</li>
                <li>Enter/Space: Activate buttons and links</li>
                <li>Escape: Close modals and menus</li>
                <li>Alt+M: Skip to main content</li>
                <li>Arrow keys: Navigate within lists and menus</li>
              </ul>
            </div>
          </div>
        </section>

        {/* Color Contrast */}
        <section aria-labelledby="contrast-heading">
          <h2 id="contrast-heading" className="text-2xl font-semibold mb-6">
            Color Contrast Tests
          </h2>
          <div className="space-y-4 p-6 border rounded-lg">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div className="space-y-2">
                <h3 className="font-medium">Normal Text Contrast</h3>
                <div className="bg-white text-black p-4 rounded">Black on White</div>
                <div className="bg-gray-100 text-gray-900 p-4 rounded">Dark Gray on Light Gray</div>
                <div className="bg-blue-500 text-white p-4 rounded">White on Blue</div>
              </div>
              <div className="space-y-2">
                <h3 className="font-medium">Large Text Contrast</h3>
                <div className="bg-white text-gray-600 p-4 rounded text-lg">
                  Medium Gray on White (Large)
                </div>
                <div className="bg-gray-800 text-gray-200 p-4 rounded text-lg">
                  Light Gray on Dark Gray (Large)
                </div>
              </div>
            </div>
          </div>
        </section>

        {/* Focus Indicators */}
        <section aria-labelledby="focus-heading">
          <h2 id="focus-heading" className="text-2xl font-semibold mb-6">
            Focus Indicators
          </h2>
          <div className="space-y-4 p-6 border rounded-lg">
            <p className="text-muted-foreground mb-4">
              Tab through these elements to see focus indicators:
            </p>
            <div className="space-y-4">
              <Button>Primary Button</Button>
              <Button variant="outline">Outline Button</Button>
              <Button variant="secondary">Secondary Button</Button>
              <Input placeholder="Input field with focus" />
              <textarea
                className="flex w-full rounded-2xl border border-input bg-background px-4 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50 resize-none"
                placeholder="Textarea with focus"
                rows={3}
              />
            </div>
          </div>
        </section>
      </main>

      <footer className="mt-16 p-6 border-t">
        <p className="text-center text-muted-foreground">
          This test page demonstrates WCAG 2.1 AA compliance features. Use your browser's
          accessibility tools and screen reader software to test these components.
        </p>
      </footer>
    </div>
  );
}
