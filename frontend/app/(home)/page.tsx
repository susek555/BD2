import SideBar from "@/app/ui/(home)/sidebar";

export default function Home() {
  return (
  <main>
    <div className="flex flex-col md:flex-row flex-grow">
      <div className="w-full md:w-80 py-4 flex-none">
      <SideBar />
      </div>
      <div className="flex-grow p-6 md:overflow-y-auto md:p-12">
      <p>Main Page</p>
      </div>
    </div>
  </main>
  );
}
