import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'

export default function AuthPage() {
  return (
    <section className="w-full py-16 md:py-24">
      <div className="container mx-auto px-4 max-w-md">
        <h1 className="text-3xl md:text-4xl font-bold tracking-tight mb-6">Sign in</h1>
        <form className="space-y-4">
          <Input type="email" placeholder="Email" />
          <Input type="password" placeholder="Password" />
          <Button className="w-full">Sign in</Button>
        </form>
        <div className="text-sm text-muted-foreground mt-4">Demo credentials: demo@blytz.app / demo123</div>
      </div>
    </section>
  )
}