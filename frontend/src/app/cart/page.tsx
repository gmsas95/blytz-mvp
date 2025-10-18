import { api } from '@/lib/api-adapter'
import { Button } from '@/components/ui/button'

export default async function CartPage() {
  const res = await api.getCart()
  const cart = res.success && res.data ? res.data : null

  return (
    <section className="w-full py-16 md:py-24">
      <div className="container mx-auto px-4 space-y-6">
        <h1 className="text-3xl md:text-4xl font-bold tracking-tight">Cart</h1>
        {!cart || cart.items.length === 0 ? (
          <div className="text-muted-foreground">Your cart is empty.</div>
        ) : (
          <div className="space-y-4">
            <ul className="space-y-3">
              {cart.items.map((item) => (
                <li key={item.id} className="flex items-center justify-between rounded-xl border p-4">
                  <div>
                    <div className="font-medium">{item.product.title}</div>
                    <div className="text-sm text-muted-foreground">Qty: {item.quantity}</div>
                  </div>
                  <div className="text-sm">${(item.product.price * item.quantity).toFixed(2)}</div>
                </li>
              ))}
            </ul>
            <div className="flex items-center justify-between pt-4 border-t">
              <div className="text-sm text-muted-foreground">Items: {cart.itemCount}</div>
              <div className="text-lg font-semibold">Total: ${cart.total.toFixed(2)}</div>
            </div>
            <a href="/checkout"><Button size="lg">Proceed to Checkout</Button></a>
          </div>
        )}
      </div>
    </section>
  )
}