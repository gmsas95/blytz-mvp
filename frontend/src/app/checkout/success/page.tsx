'use client'

import { useEffect, useState } from 'react'
import { useRouter, useSearchParams } from 'next/navigation'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { CheckCircle, Loader2, AlertCircle } from 'lucide-react'
import { api } from '@/lib/api-adapter'
import { PaymentResponse } from '@/types'

export default function CheckoutSuccessPage() {
  const router = useRouter()
  const searchParams = useSearchParams()
  const [loading, setLoading] = useState(true)
  const [payment, setPayment] = useState<PaymentResponse | null>(null)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    const paymentId = searchParams.get('payment_id')
    const orderNumber = searchParams.get('order_number')

    if (paymentId) {
      verifyPayment(paymentId)
    } else if (orderNumber) {
      // Fallback for redirect-based payments
      setPayment({
        paymentId: '',
        orderNumber,
        amount: 0,
        currency: 'MYR',
        status: 'success',
        paymentMethod: 'unknown',
        redirectUrl: '',
        createdAt: new Date().toISOString(),
        expiresAt: ''
      })
      setLoading(false)
    } else {
      setError('No payment information found')
      setLoading(false)
    }
  }, [searchParams])

  const verifyPayment = async (paymentId: string) => {
    try {
      const response = await api.getPaymentStatus(paymentId)
      
      if (response.success && response.data) {
        setPayment(response.data)
      } else {
        setError('Failed to verify payment status')
      }
    } catch (err) {
      setError('Error verifying payment')
    } finally {
      setLoading(false)
    }
  }

  if (loading) {
    return (
      <div className="container mx-auto px-4 py-8">
        <div className="flex items-center justify-center min-h-[400px]">
          <div className="text-center">
            <Loader2 className="w-8 h-8 animate-spin mx-auto mb-4" />
            <p>Verifying your payment...</p>
          </div>
        </div>
      </div>
    )
  }

  if (error) {
    return (
      <div className="container mx-auto px-4 py-8">
        <div className="max-w-md mx-auto text-center">
          <AlertCircle className="w-16 h-16 text-red-500 mx-auto mb-4" />
          <h1 className="text-3xl font-bold mb-4">Payment Verification Failed</h1>
          <p className="text-muted-foreground mb-6">{error}</p>
          <div className="space-y-3">
            <Button onClick={() => router.push('/checkout')} className="w-full">
              Try Again
            </Button>
            <Button 
              variant="outline" 
              onClick={() => router.push('/')}
              className="w-full"
            >
              Return to Home
            </Button>
          </div>
        </div>
      </div>
    )
  }

  return (
    <div className="container mx-auto px-4 py-8">
      <div className="max-w-2xl mx-auto">
        <div className="text-center mb-8">
          <CheckCircle className="w-16 h-16 text-green-500 mx-auto mb-4" />
          <h1 className="text-3xl font-bold mb-4">Payment Successful!</h1>
          <p className="text-muted-foreground">
            Thank you for your purchase. Your order has been confirmed.
          </p>
        </div>

        {payment && (
          <Card>
            <CardHeader>
              <CardTitle>Order Details</CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <p className="text-sm text-muted-foreground">Order Number</p>
                  <p className="font-medium">{payment.orderNumber}</p>
                </div>
                <div>
                  <p className="text-sm text-muted-foreground">Payment ID</p>
                  <p className="font-medium">{payment.paymentId || 'N/A'}</p>
                </div>
                <div>
                  <p className="text-sm text-muted-foreground">Amount Paid</p>
                  <p className="font-medium">RM{payment.amount.toFixed(2)}</p>
                </div>
                <div>
                  <p className="text-sm text-muted-foreground">Payment Method</p>
                  <p className="font-medium capitalize">{payment.paymentMethod}</p>
                </div>
                <div>
                  <p className="text-sm text-muted-foreground">Status</p>
                  <p className="font-medium text-green-600">Completed</p>
                </div>
                <div>
                  <p className="text-sm text-muted-foreground">Date</p>
                  <p className="font-medium">
                    {new Date(payment.createdAt).toLocaleDateString()}
                  </p>
                </div>
              </div>

              <div className="bg-blue-50 p-4 rounded-lg">
                <h3 className="font-medium mb-2">What happens next?</h3>
                <ul className="text-sm space-y-1 text-muted-foreground">
                  <li>• You'll receive an email confirmation shortly</li>
                  <li>• Order details will be available in your account</li>
                  <li>• For auction items, the seller will contact you for delivery</li>
                  <li>• You can track your order status in real-time</li>
                </ul>
              </div>
            </CardContent>
          </Card>
        )}

        <div className="mt-8 space-y-3">
          <Button onClick={() => router.push('/orders')} className="w-full">
            View Order History
          </Button>
          <Button 
            variant="outline" 
            onClick={() => router.push('/')}
            className="w-full"
          >
            Continue Shopping
          </Button>
        </div>
      </div>
    </div>
  )
}