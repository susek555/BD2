import SideBar from "@/app/ui/[main]/sidebar";
import TopBar from "@/app/ui/[main]/topbar";

export default function Home() {
  return (
  <main>
    <div className="flex flex-grow flex-col">
      <div className="h-full flex-none md:h-8">
        <TopBar />
      </div>
      <div className="flex flex-grow flex-col md:flex-row">
        <div className="w-full flex-none md:w-64">
          <SideBar />
        </div>
        <div className="flex-grow p-6 md:overflow-y-auto md:p-12">
          <p>Main Page</p>
        </div>
      </div>
    </div>
  </main>
  );
}
