package logar

var index_html = `{{block "index" .}}
<!DOCTYPE html>
<html lang="en">

<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">

  <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.7.2/css/all.min.css"
    integrity="sha512-Evv84Mr4kqVGRNSgIGL/F/aIDqQb7xQ2vcrdIwxfjThSH8CSR7PBEakCr51Ck+w+/U6swU2Im1vVX0SVk9ABhg=="
    crossorigin="anonymous" referrerpolicy="no-referrer" />
  <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/css/bootstrap.min.css" rel="stylesheet"
    integrity="sha384-QWTKZyjpPEjISv5WaRU9OFeRpok6YctnYmDr5pNlyT2bRjXh0JMhjY6hW+ALEwIH" crossorigin="anonymous">
  <script src="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/js/bootstrap.bundle.min.js"
    integrity="sha384-YvpcrYf0tY3lHB60NNkmXc5s9fDVZLESaAA55NDzOxhy9GkcIdslK1eN7N6jIeHz"
    crossorigin="anonymous"></script>
  <script src="https://unpkg.com/htmx.org@2.0.4"
    integrity="sha384-HGfztofotfshcF7+8n44JQL2oJmowVChPTg48S+jvZoztPfvwD79OC/LTtG6dMp+"
    crossorigin="anonymous"></script>

  <link rel="icon" type="image/jpeg" href="">
  <title> Logger Web </title>

  <style>
    @import url('https://fonts.googleapis.com/css2?family=Open+Sans:ital,wght@0,300..800;1,300..800&display=swap');

    * {
      margin: 0;
      padding: 0;
      box-sizing: border-box;
    }

    p,
    h1,
    h2,
    h3,
    h4,
    h5,
    h6 {
      font-family: "Open Sans", serif;
      font-optical-sizing: auto;
    }

    .fade-in.htmx-added {
      opacity: 0;
    }
    .fade-in {
      opacity: 1;
      transition: calc(var(--index) * 30ms) ease-in-out;
    }

    .table-log {
      --bs-table-bg: #d1e7dd;
      --bs-table-striped-bg: #c7dbd2;

      --bs-table-bg: #C4C4C4;
      --bs-table-striped-bg: #D6D6D6;

      --bs-table-bg: #d1e7dd;
      --bs-table-striped-bg: #c7dbd2;
    }

    .table-info {
      --bs-table-bg: #cff4fc;
      --bs-table-striped-bg: #c5e8ef;

      --bs-table-bg: #C5DAF2;
      --bs-table-striped-bg: #D8E6F8;
    }

    .table-warn {
      --bs-table-bg: #fff3cd;
      --bs-table-striped-bg: #F7E6B5;
      
      --bs-table-bg: #F7E3A3;
      --bs-table-striped-bg: #FAEDD2;

      --bs-table-bg: #fff3cd;
      --bs-table-striped-bg: #F7E6B5;
    }
    
    .table-error {
      --bs-table-bg: #f8d7da;
      --bs-table-striped-bg: #eccccf;
      
      --bs-table-bg: #F6D1D1;
      --bs-table-striped-bg: #F9E1E1;

      --bs-table-bg: #f8d7da;
      --bs-table-striped-bg: #eccccf;
    }
    
    .table-crit {
      --bs-table-bg: #d6959b;
      --bs-table-striped-bg: #c98b91;
      
      --bs-table-bg: #D98A8A;
      --bs-table-striped-bg: #EBA8A8;
    }
    
    .table-fatal {
      --bs-table-bg: #582b2f;
      --bs-table-striped-bg: #4e2428;

      --bs-table-bg: #D68C8C;
      --bs-table-striped-bg: #E6B1B1;

      font-weight: bold;
      /*
      --bs-table-color: #fff;
      --bs-table-striped-color: #fff;
      */
    }

    .table-debug {
      --bs-table-bg: #6ad4a5;
      --bs-table-striped-bg: #5ecb9a;

      --bs-table-bg: #C7D9F0;
      --bs-table-striped-bg: #D9E7F7;
    }

    .table-trace {
      --bs-table-bg: #6ad4a5;
      --bs-table-striped-bg: #5ecb9a;

      --bs-table-bg: #e6e6e6;
      --bs-table-striped-bg: #d9d9d9;
    }

    mark {
      background-color: lime;
    }
  </style>
</head>

<body>
  <div class="d-flex">
    {{template "sidebar" .}}

    <div class="px-3" style="max-height: 99vh; width: 100%; overflow-y: scroll;">
      <div class="mx-3">
        <h1 class="mt-3">{{.AppName}} logger</h1>
      </div>

      <form hx-get="/logger/{{or .CurrentRoute "all"}}/logs" hx-trigger="change, keyup changed delay:0.3s" hx-target="tbody" hx-swap="innerHTML">
        <!--
        <div class="d-flex px-5 py-3 gap-2">
          
          <button class="btn btn-outline-secondary" type="submit">Filter</button>
        </div>
      -->
        <div class="d-flex px-5 py-3 gap-2">
          <input name="filter" type="text" class="form-control" placeholder="Search" aria-label="Search"
          aria-describedby="button-addon2">
          <select class="form-select" name="severity" aria-label="Default select example">
            <option value="0">All</option>
            <option value="1">Log</option>
            <option value="2">Info</option>
            <option value="3">Warning</option>
            <option value="4">Error</option>
            <option value="5">Fatal</option>
            <option value="6">Trace</option>
          </select>

          <button class="btn btn-outline-secondary" type="submit" id="button-addon2">Search</button>
        </div>
      </form>

      
        <table class="table table-striped table-bordered text-center" style="height: 400px; overflow-y: scroll;">
          <thead>
            <tr class="table-dark">
              <th scope="col" style="width: 1%;">#</th>
              <th scope="col" style="width: 1%;"><i class="fa-solid fa-signal"></i></th>
              <th scope="col" style="width: 50ch;">Timestamp</th>
              <th scope="col" style="width: 70%;">Message</th>
              <th scope="col" style="width: 1%;">Category</th>
            </tr>
          </thead>
          <tbody>
          
            {{template "logs" .Logs}}
          </tbody>
        </table>
    </div>
  </div>
</body>

</html>
{{end}}

{{define "logs"}}
  {{$data := .}}
  {{range $index, $element := .Logs}}
  <tr
  {{ if and (ne $data.LastID 0) (eq $index (Minus (len $data.Logs) 1)) }}
    data-last="true"
    hx-get="/logger/{{or $data.Model "all"}}/logs?cursor={{$data.LastID}}"
    hx-trigger="intersect once"
    hx-swap="afterend"
  {{end}}
    style="--index: {{$index}};"
    class="fade-in table-success table-{{ .Severity | Severity_String | ToLower }}">
    <th scope="row">{{.ID}}</th>
    <td>{{.Severity | Severity_String | ToUpper}}</td>
    <td>{{.CreatedAt.Format "2006-01-02 15:04:05.000"}}</td>
    <td style="text-align: left; word-wrap: break-word; word-break: break-all;">{{Escape .Message}}</td>
    <td>{{.Category}}</td>
</tr>
  {{end}}

{{end}}

{{define "sidebar"}}
<div class="d-flex flex-column flex-shrink-0 p-3 text-white bg-dark" style="width: 280px; height: 100vh;">
  <a href="/logger/" class="d-flex align-items-center mb-3 mb-md-0 me-md-auto text-white text-decoration-none">
    &nbsp;&nbsp; <i class="fa-solid fa-truck"></i> &nbsp;&nbsp;
    <span class="fs-4">Logger</span>
  </a>
  <hr>
  <ul class="nav nav-pills flex-column mb-auto">

    {{ $data := . }}
    {{range $route := .Routes}}
    <li class="nav-item">
      <a hx-boost="true" href="{{$route.ID}}" class="nav-link text-white {{if eq $data.CurrentRoute $route.ID}}active{{end}}">
        <i class="fa-solid fa-angle-right"></i>&nbsp;
        {{$route.Name}}
      </a>
    </li>
    {{end}}
  </ul>

  <hr>
  <div class="dropdown">
    <a href="#" class="d-flex align-items-center text-white text-decoration-none dropdown-toggle" id="dropdownUser"
      data-bs-toggle="dropdown" aria-expanded="false">
      <img
        src=""
        alt="" width="32" height="32" class="rounded-circle me-2">
      <strong>User</strong>
    </a>
    <ul class="dropdown-menu dropdown-menu-dark text-small shadow" aria-labelledby="dropdownUser">
      <li><a class="dropdown-item" href="#">New project...</a></li>
      <li><a class="dropdown-item" href="#">Settings</a></li>
      <li><a class="dropdown-item" href="#">Profile</a></li>
      <li>
        <hr class="dropdown-divider">
      </li>
      <li><a class="dropdown-item" href="#">Sign out</a></li>
    </ul>
  </div>
</div>
{{end}}`
