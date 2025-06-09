import { TopBar } from '@/app/ui/(topbar)/topbar';
import { Toaster } from 'react-hot-toast';

export default function Layout({ children }: { children: React.ReactNode }) {
  return (
    <div className='flex flex-grow flex-col'>
      <div className='h-full flex-none md:h-10'>
        <TopBar />
      </div>
      <Toaster />

      {children}
    </div>
  );
}

// export default function Layout({ children }: { children: React.ReactNode }) {
//   return (
//     <>
//       {children}
//     </>
//   )
// }
