'use client';

import { useEffect, useRef, useState } from 'react';
import Headroom from 'react-headroom';
import TabBar from '@/app/ui/account/tab-bar';

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
      }
    );

    if (sentinelRef.current) {
      observer.observe(sentinelRef.current);
    }

    return () => observer.disconnect();
  }, []);

  return (
    <>
      <div ref={sentinelRef} className='h-0' />

      {/* Regular tab bar */}
      <div className={`${isSticky ? 'invisible' : ''}`}>
        <TabBar />
      </div>

      {/* Sticky tab bar */}
      {isSticky && (
        <Headroom
          style={{
            position: 'fixed',
            top: 0,
            left: 0,
            right: 0,
            zIndex: 50,
          }}
          className='headroom-sticky'
        >
          <div className='bg-white shadow-md'>
            <div className='mx-auto max-w-7xl px-4 sm:px-6 lg:px-8'>
              <TabBar />
            </div>
          </div>
        </Headroom>
      )}
    </>
  );
}
