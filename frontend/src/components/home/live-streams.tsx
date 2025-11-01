'use client';

import { Play, Users, Heart, Clock } from 'lucide-react';
import Link from 'next/link';
import { useEffect, useState } from 'react';

import { Button } from '@/components/ui/button';
import { Card, CardContent, CardFooter, CardHeader } from '@/components/ui/card';
import { api } from '@/lib/api-adapter';
import { Livestream } from '@/types';


export function LiveStreams() {
  const [livestreams, setLivestreams] = useState<Livestream[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    async function loadActiveLivestreams() {
      try {
        const response = await api.getActiveLivestreams();
        if (response.success && response.data) {
          setLivestreams(response.data);
        }
      } catch (error) {
        console.error('Failed to load live streams:', error);
      } finally {
        setLoading(false);
      }
    }

    loadActiveLivestreams();
  }, []);

  if (loading) {
    return (
      <section className="py-16 md:py-24 bg-secondary/30">
        <div className="container px-4 md:px-6">
          <div className="text-center space-y-4 mb-12">
            <div className="h-8 bg-muted rounded-lg w-64 mx-auto animate-pulse" />
            <div className="h-4 bg-muted rounded w-96 mx-auto animate-pulse" />
          </div>

          <div className="grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
            {[...Array(3)].map((_, i) => (
              <Card key={i} className="overflow-hidden">
                <div className="aspect-video bg-muted animate-pulse" />
                <CardHeader className="pb-3">
                  <div className="h-4 bg-muted rounded animate-pulse mb-2" />
                  <div className="h-3 bg-muted rounded w-2/3 animate-pulse" />
                </CardHeader>
                <CardContent className="pb-3">
                  <div className="h-4 bg-muted rounded w-1/3 animate-pulse mb-2" />
                  <div className="h-4 bg-muted rounded w-1/2 animate-pulse" />
                </CardContent>
                <CardFooter>
                  <div className="h-10 bg-muted rounded-lg animate-pulse w-full" />
                </CardFooter>
              </Card>
            ))}
          </div>
        </div>
      </section>
    );
  }

  if (livestreams.length === 0) {
    return null;
  }

  return (
    <section className="w-full py-16 md:py-24 bg-secondary/30">
      <div className="container mx-auto px-4">
        <div className="text-center space-y-4 mb-12">
          <div className="inline-flex items-center gap-2 bg-purple-100 text-purple-800 px-4 py-2 rounded-full text-sm font-medium">
            <div className="w-2 h-2 bg-purple-500 rounded-full animate-pulse" />
            Live Streams
          </div>

          <h2 className="text-3xl md:text-4xl font-bold tracking-tight text-foreground">
            Go Live with Sellers
          </h2>

          <p className="text-lg text-muted-foreground max-w-2xl mx-auto">
            Watch sellers showcase products live, ask questions, and bid in real-time
          </p>
        </div>

        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
          {livestreams.map((livestream) => (
            <Link key={livestream.id} href={`/livestream/${livestream.id}`} className="group">
              <Card className="overflow-hidden h-full border-border hover:shadow-lg transition-shadow duration-300">
                <div className="relative aspect-video bg-muted">
                  <img
                    src={livestream.thumbnail}
                    alt={livestream.title}
                    className="object-cover w-full h-full group-hover:scale-105 transition-transform duration-300"
                  />

                  <div className="absolute top-3 left-3">
                    <div className="inline-flex items-center gap-1 bg-red-500 text-white px-2 py-1 rounded-full text-xs font-medium">
                      <div className="w-2 h-2 bg-white rounded-full animate-pulse" />
                      LIVE
                    </div>
                  </div>

                  <div className="absolute top-3 right-3">
                    <div className="inline-flex items-center gap-1 bg-black/50 text-white px-2 py-1 rounded-full text-xs backdrop-blur-sm">
                      <Users className="w-3 h-3" />
                      {livestream.viewers}
                    </div>
                  </div>

                  <div className="absolute bottom-3 left-3 right-3">
                    <div className="text-white">
                      <div className="font-medium">{livestream.seller.name}</div>
                      <div className="text-xs text-white/80">{livestream.seller.storeName}</div>
                    </div>
                  </div>

                  <div className="absolute inset-0 flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity">
                    <div className="w-16 h-16 rounded-full bg-white/20 backdrop-blur-sm flex items-center justify-center">
                      <Play className="w-8 h-8 text-white fill-white ml-1" />
                    </div>
                  </div>
                </div>

                <CardHeader className="pb-3">
                  <h3 className="font-semibold text-lg leading-tight group-hover:text-primary transition-colors">
                    {livestream.title}
                  </h3>
                  <p className="text-sm text-muted-foreground line-clamp-2">
                    {livestream.description}
                  </p>
                </CardHeader>

                <CardContent className="pb-3">
                  <div className="flex items-center justify-between text-sm">
                    <div className="flex items-center gap-1 text-muted-foreground">
                      <Heart className="w-4 h-4" />
                      <span>{livestream.likes}</span>
                    </div>

                    <div className="flex items-center gap-1 text-muted-foreground">
                      <Clock className="w-4 h-4" />
                      <span>{Math.floor(livestream.duration / 60)}m</span>
                    </div>

                    <div className="flex items-center gap-1 text-muted-foreground">
                      <Users className="w-4 h-4" />
                      <span>{livestream.products.length} items</span>
                    </div>
                  </div>
                </CardContent>

                <CardFooter>
                  <Button className="w-full gap-2">
                    <Play className="w-4 h-4" />
                    Watch Live
                  </Button>
                </CardFooter>
              </Card>
            </Link>
          ))}
        </div>

        <div className="text-center mt-12">
          <Link href="/livestream">
            <Button size="lg" variant="outline">
              Browse All Streams
            </Button>
          </Link>
        </div>
      </div>
    </section>
  );
}
