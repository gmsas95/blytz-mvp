import { notFound } from 'next/navigation'
import { api } from '@/lib/api-adapter'
import { Button } from '@/components/ui/button'

interface Props { params: { id: string } }

export default async function LivestreamDetailPage({ params }: Props) {
  const res = await api.getLivestream(params.id)
  if (!res.success || !res.data) return notFound()
  const ls = res.data

  return (
    <section className="w-full py-16 md:py-24">
      <div className="container mx-auto px-4 grid gap-10 lg:grid-cols-2 items-start">
        <div className="w-full">
          <div className="aspect-video bg-muted rounded-2xl overflow-hidden">
            {/* eslint-disable-next-line @next/next/no-img-element */}
            <img src={ls.thumbnail} alt={ls.title} className="object-cover w-full h-full" />
          </div>
        </div>
        <div className="space-y-6">
          <h1 className="text-3xl md:text-4xl font-bold tracking-tight">{ls.title}</h1>
          <p className="text-muted-foreground">{ls.description}</p>
          <div className="flex gap-3">
            <Button size="lg">Watch Live</Button>
            <a href="/products"><Button size="lg" variant="outline">Browse Items</Button></a>
          </div>
        </div>
      </div>
    </section>
  )
}