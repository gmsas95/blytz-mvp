import Link from 'next/link';

import { Button } from '@/components/ui/button';
import { Card, CardFooter, CardHeader } from '@/components/ui/card';
import { api } from '@/lib/api-adapter';

export default async function LivestreamListPage() {
  const res = await api.getLivestreams();
  const items = res.success && res.data ? res.data.items : [];

  return (
    <section className="w-full py-16 md:py-24 bg-secondary/30">
      <div className="container mx-auto px-4">
        <div className="text-center space-y-4 mb-12">
          <h1 className="text-3xl md:text-4xl font-bold tracking-tight">Live Streams</h1>
          <p className="text-muted-foreground">Watch sellers showcase products live</p>
        </div>
        {items.length === 0 ? (
          <div className="text-center text-muted-foreground">No streams found.</div>
        ) : (
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
            {items.map((ls) => (
              <Link key={ls.id} href={`/livestream/${ls.id}`} className="group">
                <Card className="overflow-hidden h-full">
                  <div className="relative aspect-video bg-muted">
                    {/* eslint-disable-next-line @next/next/no-img-element */}
                    <img
                      src={ls.thumbnail}
                      alt={ls.title}
                      className="object-cover w-full h-full group-hover:scale-105 transition-transform duration-300"
                    />
                  </div>
                  <CardHeader className="pb-3">
                    <h3 className="font-semibold text-lg leading-tight group-hover:text-primary transition-colors">
                      {ls.title}
                    </h3>
                    <p className="text-sm text-muted-foreground line-clamp-2">{ls.description}</p>
                  </CardHeader>
                  <CardFooter>
                    <Button className="w-full">Watch</Button>
                  </CardFooter>
                </Card>
              </Link>
            ))}
          </div>
        )}
      </div>
    </section>
  );
}
