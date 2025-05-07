import SideBar from "@/app/ui/(home)/sidebar";
import { Suspense } from "react";


export default async function Home() {
  return (
  <main>
    <div className="flex flex-col md:flex-row flex-grow">
      <div className="w-full md:w-80 py-4 h-full flex-none">
        <Suspense>
          <SideBar />
        </Suspense>
      </div>
      <div className="flex-grow p-6 md:overflow-y-auto md:p-12">
        <Suspense>
          <p>Main Page</p>
        </Suspense>
      </div>
    </div>
  </main>
  );
}
