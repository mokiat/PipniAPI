<script>
  // import "App.css";
  import Tailwind from "./Tailwind.svelte";

  import {
    Navbar,
    NavBrand,
    NavLi,
    NavUl,
    NavHamburger,
    Dropdown,
    DropdownItem,
    DropdownDivider,
  } from "flowbite-svelte";
  import { Icon } from "flowbite-svelte-icons";
</script>

<Tailwind />
<Navbar>
  <NavBrand>
    <img src="icon.png" class="mr-3 h-6 sm:h-9" alt="Logo" />
    <span
      class="self-center whitespace-nowrap text-xl font-semibold dark:text-white"
      >PipniAPI</span
    >
  </NavBrand>
  <NavUl>
    <NavLi id="nav-menu1" class="cursor-pointer">
      Production<Icon
        name="chevron-down-outline"
        class="w-3 h-3 ml-2 text-primary-800 dark:text-white inline"
      />
    </NavLi>
    <Dropdown triggeredBy="#nav-menu1" class="w-44 z-20">
      <DropdownItem>Production</DropdownItem>
      <DropdownItem>Staging</DropdownItem>
      <DropdownDivider />
      <DropdownItem>Settings</DropdownItem>
    </Dropdown>
  </NavUl>
</Navbar>

<h1>Hello Svelte</h1>

<style>
  h1 {
    color: red;
    font-style: italic;
  }
</style>
