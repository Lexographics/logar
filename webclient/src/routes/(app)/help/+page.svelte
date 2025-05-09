<script>
  import BaseView from "$lib/widgets/BaseView.svelte";
  import LL from "../../../i18n/i18n-svelte";
</script>

<BaseView>
  <div class="help-container">
    <h2>{$LL.help.header()}</h2>
    
    <section>
      <h3>{$LL.help.about_logar.title()}</h3>
      <p>{$LL.help.about_logar.content()}</p>
    </section>
    
    <section>
      <h3>{$LL.help.integration_guide.title()}</h3>
      <p style="margin-bottom: 10px;">{$LL.help.integration_guide.content()}</p>

      <a href="https://godoc.org/github.com/Lexographics/logar" target="_blank">
        <img src="https://pkg.go.dev/badge/github.com/Lexographics/logar.svg" alt="Go Reference">
      </a>
    </section>
    
    <section>
      <h3>{$LL.help.troubleshooting.title()}</h3>
      <p>
        {$LL.help.troubleshooting.pre_twitter()} <a href="https://x.com/ibrsdk" target="_blank">@ibrsdk</a>
        {$LL.help.troubleshooting.pre_github()} <a href="https://github.com/Lexographics/logar" target="_blank">{$LL.help.troubleshooting.github_text()}</a>.
      </p>
    </section>
    
    <section>
      <h3>{$LL.help.additional_resources.title()}</h3>
      <ul>
        <li><a href="https://github.com/Lexographics/logar" target="_blank">{$LL.help.additional_resources.github_repository()}</a></li>
        <li><a href="https://godoc.org/github.com/Lexographics/logar" target="_blank">{$LL.help.additional_resources.api_docs()}</a></li>
      </ul>
    </section>
  </div>
</BaseView>

<style>
  .help-container {
    max-width: 800px;
    margin: 0 auto;
    padding: 20px 50px;
    height: 100%;
    background-color: var(--card-background);
  }
  
  section {
    margin-bottom: 30px;
  }
  
  h2 {
    font-size: 2rem;
    margin-bottom: 20px;
  }
  
  h3 {
    font-size: 1.5rem;
    margin: 20px 0 10px;
    border-bottom: 1px solid #eee;
    padding-bottom: 5px;
  }
  
  li {
    margin-bottom: 5px;
  }
</style>