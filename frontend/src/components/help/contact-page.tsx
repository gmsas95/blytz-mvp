'use client';

import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';

export function ContactPage() {
  return (
    <div className="container mx-auto px-4 py-8 max-w-4xl">
      <div className="mb-8">
        <h1 className="text-4xl font-bold mb-4">Contact Us</h1>
        <p className="text-lg text-muted-foreground">
          Get in touch with our support team for help with your account, orders, or any questions.
        </p>
      </div>

      <div className="grid lg:grid-cols-2 gap-8 mb-8">
        <Card>
          <CardHeader>
            <CardTitle>Send us a Message</CardTitle>
          </CardHeader>
          <CardContent>
            <form className="space-y-4">
              <div>
                <label className="block text-sm font-medium mb-2">Email</label>
                <input 
                  type="email" 
                  className="w-full px-3 py-2 border rounded-md"
                  placeholder="your@email.com"
                />
              </div>
              <div>
                <label className="block text-sm font-medium mb-2">Subject</label>
                <select className="w-full px-3 py-2 border rounded-md">
                  <option>General Question</option>
                  <option>Order Issue</option>
                  <option>Account Help</option>
                  <option>Technical Support</option>
                  <option>Partnership Inquiry</option>
                </select>
              </div>
              <div>
                <label className="block text-sm font-medium mb-2">Message</label>
                <textarea 
                  className="w-full px-3 py-2 border rounded-md h-32"
                  placeholder="How can we help you?"
                />
              </div>
              <Button className="w-full">Send Message</Button>
            </form>
          </CardContent>
        </Card>

        <div className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle>Other Ways to Reach Us</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                <div>
                  <h4 className="font-semibold">Email Support</h4>
                  <p className="text-muted-foreground">support@blytz.app</p>
                  <p className="text-sm">Response within 24 hours</p>
                </div>
                <div>
                  <h4 className="font-semibold">Phone Support</h4>
                  <p className="text-muted-foreground">+1 (555) 123-4567</p>
                  <p className="text-sm">Mon-Fri, 9am-6pm EST</p>
                </div>
                <div>
                  <h4 className="font-semibold">Live Chat</h4>
                  <p className="text-muted-foreground">Available on website</p>
                  <p className="text-sm">Mon-Fri, 9am-9pm EST</p>
                </div>
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle>Frequently Asked Questions</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="space-y-3">
                <div>
                  <h4 className="font-medium">How do I track my order?</h4>
                  <p className="text-sm text-muted-foreground">
                    Check your order history for tracking information.
                  </p>
                </div>
                <div>
                  <h4 className="font-medium">What is your return policy?</h4>
                  <p className="text-sm text-muted-foreground">
                    30-day return window for most items.
                  </p>
                </div>
                <div>
                  <h4 className="font-medium">How do I become a seller?</h4>
                  <p className="text-sm text-muted-foreground">
                    Create a seller account and complete verification.
                  </p>
                </div>
              </div>
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  );
}