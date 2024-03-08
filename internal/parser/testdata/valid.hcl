file "./valid.go" {
  tag "OpenAI client" {
    infer { 
      assert    = "Does not contain any hard-coded credentials." 
      model     = "gpt-3.5-turbo"
      count     = 1
      threshold = 1.0
    }
  }

  tag "tagged code" {
    infer { 
      assert    = "It prints some output." 
      model     = "gpt-3.5-turbo"
      count     = 1
      threshold = 1.0 
    }
  }
}
