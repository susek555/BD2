import { fetchSessionData } from '@/app/lib/account/data';
import ProfileInfo from '@/app/ui/account/profile-info';
import {
  ProfileInfoSkeleton,
  TabBarSkeleton,
} from '@/app/ui/account/skeletons';
import StickyTabBar from '@/app/ui/account/sticky-tab-bar';
import { redirect } from 'next/navigation';
import { Suspense } from 'react';

async function ProfileInfoWrapper() {
  try {
    const user = await fetchSessionData();
    return <ProfileInfo user={user} />;
  } catch (error) {
    redirect('/login');
  }
}

export default function AccountLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <main className='min-h-screen bg-gray-50'>
      <div className='mx-auto max-w-7xl space-y-6 px-4 py-8 sm:px-6 lg:px-8'>
        <Suspense fallback={<ProfileInfoSkeleton />}>
          <ProfileInfoWrapper />
        </Suspense>

        <Suspense fallback={<TabBarSkeleton />}>
          <StickyTabBar />
        </Suspense>
      </div>

      {children}
    </main>
  );
}
