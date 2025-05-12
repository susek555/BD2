'use client';

import { profileTabs, Tab } from '@/app/lib/definitions';
import clsx from 'clsx';
import Link from 'next/link';
import { usePathname } from 'next/navigation';

interface ProfileTabsProps {
  tabs: Tab[];
}

export function ProfileTabs({ tabs }: ProfileTabsProps) {
  const pathname = usePathname();
  const activeTab = pathname.split('/').pop() || 'activity';

  return (
    <div className='border-b border-gray-200 bg-white'>
      <nav className='flex gap-2 px-4 md:px-6' aria-label='Profile tabs'>
        {tabs.map((tab) => (
          <Link
            key={tab.id}
            href={`/account/${tab.id}`}
            className={clsx(
              'border-b-2 px-4 py-2 text-sm font-medium transition-colors',
              activeTab === tab.id
                ? 'border-blue-500 text-blue-600'
                : 'border-transparent text-gray-500 hover:border-gray-300 hover:text-gray-700',
            )}
            aria-current={activeTab === tab.id ? 'page' : undefined}
          >
            {tab.label}
          </Link>
        ))}
      </nav>
    </div>
  );
}

export default function TabBar() {
  return <ProfileTabs tabs={profileTabs} />;
}
