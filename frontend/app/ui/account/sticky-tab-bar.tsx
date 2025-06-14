'use client';

import TabBar from '@/app/ui/account/tab-bar';
import { useEffect, useRef, useState } from 'react';
import Headroom from 'react-headroom';

export default function StickyTabBar() {
  const [isSticky, setIsSticky] = useState(false);
  const sentinelRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const observer = new IntersectionObserver(
      ([entry]) => {
        setIsSticky(!entry.isIntersecting);
      },
      {
        rootMargin: '-1px 0px 0px 0px',
        threshold: 0,
      },
    );

    if (sentinelRef.current) {
      observer.observe(sentinelRef.current);
    }

    return () => observer.disconnect();
  }, []);

  return (
    <>
      <div ref={sentinelRef} className='h-0' />

      <div>
        <TabBar />
      </div>

      <div
        style={{
          position: 'fixed',
          top: 0,
          left: 0,
          right: 0,
          zIndex: 50,
          height: isSticky ? 'auto' : 0,
          overflow: 'hidden',
        }}
      >
        <Headroom disable={!isSticky}>
          <div className='bg-white shadow-md'>
            <div className='mx-auto max-w-7xl px-4 sm:px-6 lg:px-8'>
              <TabBar />
            </div>
          </div>
        </Headroom>
      </div>
    </>
  );
}
