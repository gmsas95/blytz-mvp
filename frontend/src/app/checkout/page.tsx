'use client'

import { useState, useEffect } from 'react'
import { useRouter } from 'next/navigation'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { RadioGroup, RadioGroupItem } from '@/components/ui/radio-group'
import { Label } from '@/components/ui/label'
import { Separator } from '@/components/ui/separator'
import { Badge } from '@/components/ui/badge'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { Loader2, CreditCard, Smartphone, Building2, CheckCircle } from 'lucide-react'
import { api } from '@/lib/api-adapter'
import { Cart, PaymentMethodInfo, PaymentRequest, PaymentResponse } from '@/types'

declare global {
  interface Window {
    IPGSeamless: any
  }
}

export default function CheckoutPage() {
  const router = useRouter()
  const [cart, setCart] = useState<Cart | null>(null)
  const [paymentMethods, setPaymentMethods] = useState<PaymentMethodInfo[]>([])
  const [selectedMethod, setSelectedMethod] = useState<string>('fpx')
  const [loading, setLoading] = useState(true)
  const [processing, setProcessing] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [fiuuConfig, setFiuuConfig] = useState<any>(null)
  const [paymentComplete, setPaymentComplete] = useState(false)
  const [paymentResponse, setPaymentResponse] = useState<PaymentResponse | null>(null)

  useEffect(() => {
    loadCheckoutData()
  }, [])

  const loadCheckoutData = async () => {
    try {
      setLoading(true)
      const [cartResponse, methodsResponse, configResponse] = await Promise.all([
        api.getCart(),
        api.getPaymentMethods(),
        api.getFiuuSeamlessConfig()
      ])

      if (cartResponse.success && cartResponse.data) {
        setCart(cartResponse.data)
      }

      if (methodsResponse.success && methodsResponse.data) {
        setPaymentMethods(methodsResponse.data)
      }

      if (configResponse.success && configResponse.data) {
        setFiuuConfig(configResponse.data)
      }
    } catch (err) {
      setError('Failed to load checkout data')
    } finally {
      setLoading(false)
    }
  }

  const getPaymentIcon = (method: PaymentMethodInfo) => {
    switch (method.type) {
      case 'bank_transfer':
        return <Building2 className="w-5 h-5" />
      case 'ewallet':
        return <Smartphone className="w-5 h-5" />
      case 'card':
        return <CreditCard className="w-5 h-5" />
      default:
        return <CreditCard className="w-5 h-5" />
    }
  }

  const calculateTotal = () => {
    if (!cart) return 0
    const selectedPaymentMethod = paymentMethods.find(m => m.id === selectedMethod)
    const processingFee = selectedPaymentMethod?.fee || 0
    return cart.total + processingFee
  }

  const handlePayment = async () => {
    if (!cart || !fiuuConfig) return

    try {
      setProcessing(true)
      setError(null)

      // Update Fiuu config with current cart data
      const updatedConfig = {
        ...fiuuConfig,
        amount: calculateTotal(),
        orderNumber: `BLYTZ_${Date.now()}`,
        productDescription: `Payment for ${cart.itemCount} item(s)`,
        paymentMethod: selectedMethod === 'fpx' ? 'all' : selectedMethod
      }

      // Create payment
      const paymentRequest: PaymentRequest = {
        amount: calculateTotal(),
        currency: 'MYR',
        paymentMethod: selectedMethod,
        orderNumber: updatedConfig.orderNumber,
        description: updatedConfig.productDescription,
        returnUrl: `${window.location.origin}/checkout/success`,
        cancelUrl: `${window.location.origin}/checkout/cancel`,
        webhookUrl: `${window.location.origin}/api/payments/webhook`
      }

      const paymentResult = await api.createPayment(paymentRequest)
      
      if (!paymentResult.success) {
        throw new Error(paymentResult.error || 'Failed to create payment')
      }

      if (paymentResult.data) {
        setPaymentResponse(paymentResult.data)
      }

      // Initialize Fiuu seamless payment
      if (window.IPGSeamless) {
        const seamless = new window.IPGSeamless(updatedConfig)
        
        seamless.setCompleteCallback((response: any) => {
          console.log('Payment complete:', response)
          setPaymentComplete(true)
          setProcessing(false)
        })

        seamless.setErrorCallback((error: any) => {
          console.error('Payment error:', error)
          setError('Payment failed. Please try again.')
          setProcessing(false)
        })

        // Start payment
        seamless.makePayment()
      } else if (paymentResult.data) {
        // Fallback to redirect method
        window.location.href = paymentResult.data.redirectUrl
      }

    } catch (err) {
      setError(err instanceof Error ? err.message : 'Payment failed')
      setProcessing(false)
    }
  }

  if (loading) {
    return (
      <div className="container mx-auto px-4 py-8">
        <div className="flex items-center justify-center min-h-[400px]">
          <Loader2 className="w-8 h-8 animate-spin" />
        </div>
      </div>
    )
  }

  if (!cart || cart.items.length === 0) {
    return (
      <div className="container mx-auto px-4 py-8">
        <div className="text-center">
          <h1 className="text-3xl font-bold mb-4">Your cart is empty</h1>
          <Button onClick={() => router.push('/products')}>
            Continue Shopping
          </Button>
        </div>
      </div>
    )
  }

  if (paymentComplete && paymentResponse) {
    return (
      <div className="container mx-auto px-4 py-8">
        <div className="max-w-md mx-auto text-center">
          <CheckCircle className="w-16 h-16 text-green-500 mx-auto mb-4" />
          <h1 className="text-3xl font-bold mb-4">Payment Successful!</h1>
          <p className="text-muted-foreground mb-6">
            Your order #{paymentResponse.orderNumber} has been confirmed.
          </p>
          <div className="space-y-3">
            <Button onClick={() => router.push('/orders')} className="w-full">
              View Order
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

  const selectedPaymentMethod = paymentMethods.find(m => m.id === selectedMethod)
  const processingFee = selectedPaymentMethod?.fee || 0
  const total = calculateTotal()

  return (
    <div className="container mx-auto px-4 py-8">
      <div className="grid lg:grid-cols-2 gap-8">
        {/* Order Summary */}
        <div>
          <h2 className="text-2xl font-bold mb-6">Order Summary</h2>
          <Card>
            <CardHeader>
              <CardTitle>Cart Items ({cart.itemCount})</CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              {cart.items.map((item) => (
                <div key={item.id} className="flex items-center space-x-4">
                  <div className="w-16 h-16 bg-gray-200 rounded-lg flex items-center justify-center">
                    <span className="text-2xl">📦</span>
                  </div>
                  <div className="flex-1">
                    <h3 className="font-medium">{item.product.title}</h3>
                    <p className="text-sm text-muted-foreground">
                      Quantity: {item.quantity} × RM{item.product.price.toFixed(2)}
                    </p>
                    {item.selectedAuction && (
                      <Badge variant="secondary" className="mt-1">
                        Auction: {item.selectedAuction.product.title}
                      </Badge>
                    )}
                  </div>
                  <div className="text-right">
                    <p className="font-medium">
                      RM{(item.product.price * item.quantity).toFixed(2)}
                    </p>
                  </div>
                </div>
              ))}
              
              <Separator />
              
              <div className="space-y-2">
                <div className="flex justify-between">
                  <span>Subtotal</span>
                  <span>RM{cart.total.toFixed(2)}</span>
                </div>
                {processingFee > 0 && (
                  <div className="flex justify-between">
                    <span>Processing Fee</span>
                    <span>RM{processingFee.toFixed(2)}</span>
                  </div>
                )}
                <Separator />
                <div className="flex justify-between font-bold text-lg">
                  <span>Total</span>
                  <span>RM{total.toFixed(2)}</span>
                </div>
              </div>
            </CardContent>
          </Card>
        </div>

        {/* Payment Method */}
        <div>
          <h2 className="text-2xl font-bold mb-6">Payment Method</h2>
          
          {error && (
            <Alert variant="destructive" className="mb-6">
              <AlertDescription>{error}</AlertDescription>
            </Alert>
          )}

          <Card>
            <CardHeader>
              <CardTitle>Select Payment Method</CardTitle>
            </CardHeader>
            <CardContent>
              <RadioGroup value={selectedMethod} onValueChange={setSelectedMethod}>
                {paymentMethods.map((method) => (
                  <div key={method.id} className="flex items-center space-x-3 p-3 border rounded-lg">
                    <RadioGroupItem value={method.id} id={method.id} />
                    <Label htmlFor={method.id} className="flex items-center space-x-3 flex-1 cursor-pointer">
                      {getPaymentIcon(method)}
                      <div className="flex-1">
                        <div className="font-medium">{method.name}</div>
                        <div className="text-sm text-muted-foreground">{method.description}</div>
                      </div>
                      {method.fee > 0 && (
                        <Badge variant="secondary">
                          +RM{method.fee.toFixed(2)}
                        </Badge>
                      )}
                    </Label>
                  </div>
                ))}
              </RadioGroup>

              <Separator className="my-6" />

              <div className="space-y-4">
                <div className="bg-gray-50 p-4 rounded-lg">
                  <h3 className="font-medium mb-2">Payment Details</h3>
                  <div className="space-y-1 text-sm">
                    <div className="flex justify-between">
                      <span>Items:</span>
                      <span>{cart.itemCount}</span>
                    </div>
                    <div className="flex justify-between">
                      <span>Subtotal:</span>
                      <span>RM{cart.total.toFixed(2)}</span>
                    </div>
                    {processingFee > 0 && (
                      <div className="flex justify-between">
                        <span>Processing Fee:</span>
                        <span>RM{processingFee.toFixed(2)}</span>
                      </div>
                    )}
                    <Separator className="my-2" />
                    <div className="flex justify-between font-bold">
                      <span>Total Amount:</span>
                      <span>RM{total.toFixed(2)}</span>
                    </div>
                  </div>
                </div>

                <Button 
                  onClick={handlePayment}
                  disabled={processing}
                  className="w-full"
                  size="lg"
                >
                  {processing ? (
                    <>
                      <Loader2 className="w-4 h-4 mr-2 animate-spin" />
                      Processing Payment...
                    </>
                  ) : (
                    `Pay RM${total.toFixed(2)}`
                  )}
                </Button>

                <p className="text-xs text-muted-foreground text-center">
                  By completing this payment, you agree to our terms of service and privacy policy.
                  Your payment is secured by Fiuu payment gateway.
                </p>
              </div>
            </CardContent>
          </Card>
        </div>
      </div>

      {/* Fiuu Script */}
      {fiuuConfig && (
        <script
          src="https://sandbox.merchant.razer.com/RMS2/IPGSeamless/IPGSeamless.js"
          async
        />
      )}
    </div>
  )
}